// Copyright 2012, 2013 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package charmrevision

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/juju/clock"
	"github.com/juju/errors"
	"github.com/juju/names/v5"
	"github.com/juju/retry"
	"github.com/juju/worker/v4"
	"github.com/juju/worker/v4/catacomb"

	"github.com/juju/juju/core/arch"
	corecharm "github.com/juju/juju/core/charm"
	charmmetrics "github.com/juju/juju/core/charm/metrics"
	corehttp "github.com/juju/juju/core/http"
	"github.com/juju/juju/core/logger"
	coremodel "github.com/juju/juju/core/model"
	"github.com/juju/juju/core/os/ostype"
	"github.com/juju/juju/core/version"
	"github.com/juju/juju/core/watcher"
	"github.com/juju/juju/domain/application"
	"github.com/juju/juju/domain/application/architecture"
	applicationcharm "github.com/juju/juju/domain/application/charm"
	"github.com/juju/juju/environs/config"
	"github.com/juju/juju/internal/charm"
	internalcharm "github.com/juju/juju/internal/charm"
	"github.com/juju/juju/internal/charm/repository"
	"github.com/juju/juju/internal/charm/resource"
	"github.com/juju/juju/internal/charmhub"
	"github.com/juju/juju/internal/charmhub/transport"
	internalerrors "github.com/juju/juju/internal/errors"
)

const (
	// ErrFailedToSendMetrics is the error returned when sending metrics to the
	// charmhub fails.
	ErrFailedToSendMetrics = internalerrors.ConstError("sending metrics failed")
)

// ModelConfigService provides access to the model configuration.
type ModelConfigService interface {
	// ModelConfig returns the current config for the model.
	ModelConfig(context.Context) (*config.Config, error)

	// Watch returns a watcher that notifies of changes to the model config.
	Watch() (watcher.StringsWatcher, error)
}

// ApplicationService provides access to applications.
type ApplicationService interface {
	// GetApplicationsForRevisionUpdater returns the applications that should be
	// used by the revision updater.
	GetApplicationsForRevisionUpdater(context.Context) ([]application.RevisionUpdaterApplication, error)
}

// ModelService provides access to the model.
type ModelService interface {
	// GetModelMetrics returns the model metrics information set in the
	// database.
	GetModelMetrics(context.Context) (coremodel.ModelMetrics, error)
}

// Config defines the operation of a charm revision updater worker.
type Config struct {
	// ModelConfigService is the service used to access model configuration.
	ModelConfigService ModelConfigService

	// ApplicationService is the service used to access applications.
	ApplicationService ApplicationService

	// ModelService is the service used to access the model.
	ModelService ModelService

	// ModelTag is the tag of the model the worker is running in.
	ModelTag names.ModelTag

	// HTTPClientGetter is the getter used to create HTTP clients.
	HTTPClientGetter corehttp.HTTPClientGetter

	// NewHTTPClient is the function used to create a new HTTP client.
	NewHTTPClient NewHTTPClientFunc

	// NewCharmhubClient is the function used to create a new CharmhubClient.
	NewCharmhubClient NewCharmhubClientFunc

	// Clock is the worker's view of time.
	Clock clock.Clock

	// Period is the time between charm revision updates.
	Period time.Duration

	// Logger is the logger used for debug logging in this worker.
	Logger logger.Logger
}

// Validate returns an error if the configuration cannot be expected
// to start a functional worker.
func (config Config) Validate() error {
	if config.ModelConfigService == nil {
		return errors.NotValidf("nil ModelConfigService")
	}
	if config.ApplicationService == nil {
		return errors.NotValidf("nil ApplicationService")
	}
	if config.ModelService == nil {
		return errors.NotValidf("nil ModelService")
	}
	if config.HTTPClientGetter == nil {
		return errors.NotValidf("nil HTTPClientGetter")
	}
	if config.NewHTTPClient == nil {
		return errors.NotValidf("nil NewHTTPClient")
	}
	if config.NewCharmhubClient == nil {
		return errors.NotValidf("nil NewCharmhubClient")
	}
	if config.Clock == nil {
		return errors.NotValidf("nil Clock")
	}
	if config.Period <= 0 {
		return errors.NotValidf("non-positive Period")
	}
	if config.Logger == nil {
		return errors.NotValidf("nil Logger")
	}
	return nil
}

type revisionUpdateWorker struct {
	catacomb catacomb.Catacomb
	config   Config
}

// NewWorker returns a worker that calls UpdateLatestRevisions on the
// configured RevisionUpdater, once when started and subsequently every
// Period.
func NewWorker(config Config) (worker.Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, internalerrors.Capture(err)
	}
	w := &revisionUpdateWorker{
		config: config,
	}

	if err := catacomb.Invoke(catacomb.Plan{
		Site: &w.catacomb,
		Work: w.loop,
	}); err != nil {
		return nil, internalerrors.Capture(err)
	}

	w.config.Logger.Debugf("worker created with period %v", w.config.Period)

	return w, nil
}

// Kill is part of the worker.Worker interface.
func (w *revisionUpdateWorker) Kill() {
	w.catacomb.Kill(nil)
}

// Wait is part of the worker.Worker interface.
func (w *revisionUpdateWorker) Wait() error {
	return w.catacomb.Wait()
}

func (w *revisionUpdateWorker) loop() error {
	ctx, cancel := w.scopedContext()
	defer cancel()

	// Watch the model config for new charmhub URL values, so we can swap the
	// charmhub client to use the new URL.

	modelConfigService := w.config.ModelConfigService
	configWatcher, err := modelConfigService.Watch()
	if err != nil {
		return internalerrors.Capture(err)
	}

	if err := w.catacomb.Add(configWatcher); err != nil {
		return internalerrors.Capture(err)
	}

	logger := w.config.Logger
	logger.Debugf("watching model config for changes to charmhub URL")

	charmhubClient, err := w.getCharmhubClient(ctx)
	if err != nil {
		return internalerrors.Capture(err)
	}

	for {
		select {
		case <-w.catacomb.Dying():
			return w.catacomb.ErrDying()

		case <-w.config.Clock.After(jitter(w.config.Period)):
			w.config.Logger.Debugf("%v elapsed, performing work", w.config.Period)

			// This worker is responsible for updating the latest revision of
			// applications in the model. It does this by fetching the latest
			// revision from the charmhub and updating the model with the
			// information.
			// If the update fails, the worker will log an error and continue
			// to the next application.
			latestInfo, err := w.fetch(ctx, charmhubClient)
			if errors.Is(err, ErrFailedToSendMetrics) {
				logger.Warningf("failed to send metrics: %v", err)
				continue
			} else if err != nil {
				logger.Errorf("failed to fetch revisions: %v", err)
				continue
			}

			logger.Debugf("revisions fetched for %d applications", len(latestInfo))

			// TODO (stickupkid): Insert charms with the latest revisions.

		case change, ok := <-configWatcher.Changes():
			if !ok {
				return errors.New("model config watcher closed")
			}

			var refresh bool
			for _, key := range change {
				if key == config.CharmHubURLKey {
					refresh = true
					break
				}
			}

			if !refresh {
				continue
			}

			logger.Debugf("refreshing charmhubClient due to model config change")

			charmhubClient, err = w.getCharmhubClient(ctx)
			if err != nil {
				return internalerrors.Capture(err)
			}
		}
	}
}

func (w *revisionUpdateWorker) fetch(ctx context.Context, client CharmhubClient) ([]latestCharmInfo, error) {
	service := w.config.ApplicationService
	applications, err := service.GetApplicationsForRevisionUpdater(ctx)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	cfg, err := w.config.ModelConfigService.ModelConfig(ctx)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	buildTelemetry := cfg.Telemetry()

	if len(applications) == 0 {
		return nil, w.sendEmptyModelMetrics(ctx, client, buildTelemetry)
	}

	charmhubIDs := make([]charmhubID, len(applications))
	charmhubApps := make([]appInfo, len(applications))

	for i, app := range applications {
		charmURL, err := encodeCharmURL(app)
		if err != nil {
			w.config.Logger.Infof("encoding charm URL for %q: %v", app.Name, err)
			continue
		}

		charmhubID, err := encodeCharmhubID(app, w.config.ModelTag)
		if err != nil {
			w.config.Logger.Infof("encoding charmhub ID for %q: %v", app.Name, err)
			continue
		}

		if buildTelemetry {
			charmhubID.metrics = map[charmmetrics.MetricValueKey]string{
				charmmetrics.NumUnits:  strconv.Itoa(app.NumUnits),
				charmmetrics.Relations: strings.Join(app.Relations, ","),
			}
		}

		charmhubIDs[i] = charmhubID
		charmhubApps[i] = appInfo{
			id:       app.Name,
			charmURL: charmURL,
		}
	}

	return w.fetchInfo(ctx, client, buildTelemetry, charmhubIDs, charmhubApps)
}

func (w *revisionUpdateWorker) sendEmptyModelMetrics(ctx context.Context, client CharmhubClient, buildTelemetry bool) error {
	metadata, err := w.buildMetricsMetadata(ctx, buildTelemetry)
	if err != nil {
		return internalerrors.Capture(err)
	} else if len(metadata) == 0 {
		return nil
	}

	// Override the context which will use a shorter timeout for sending
	// metrics.
	ctx, cancel := context.WithTimeout(ctx, charmhub.RefreshTimeout)
	defer cancel()

	if err := client.RefreshWithMetricsOnly(ctx, metadata); err != nil {
		return internalerrors.Errorf("%w: %w", ErrFailedToSendMetrics, err)
	}

	return nil
}

func (w *revisionUpdateWorker) fetchInfo(ctx context.Context, client CharmhubClient, buildTelemetry bool, ids []charmhubID, apps []appInfo) ([]latestCharmInfo, error) {
	metrics, err := w.buildMetricsMetadata(ctx, buildTelemetry)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	// Override the context which will use a shorter timeout for sending
	// metrics.
	ctx, cancel := context.WithTimeout(ctx, charmhub.RefreshTimeout)
	defer cancel()

	response, err := w.request(ctx, client, metrics, ids)
	if err != nil {
		return nil, internalerrors.Errorf("requesting latest information: %w", err)
	}

	if len(response) != len(apps) {
		return nil, internalerrors.Errorf("expected %d responses, got %d", len(apps), len(response))
	}

	var latest []latestCharmInfo
	for i, result := range response {
		latest = append(latest, latestCharmInfo{
			url:       apps[i].charmURL,
			timestamp: result.timestamp,
			revision:  result.revision,
			resources: result.resources,
			appID:     apps[i].id,
		})
	}

	return latest, nil
}

// charmhubLatestCharmInfo fetches the latest information about the given
// charms from charmhub's "charm_refresh" API.
func (w *revisionUpdateWorker) request(ctx context.Context, client CharmhubClient, metrics charmhub.Metrics, ids []charmhubID) ([]charmhubResult, error) {
	configs := make([]charmhub.RefreshConfig, len(ids))
	for i, id := range ids {
		base := charmhub.RefreshBase{
			Architecture: id.arch,
			Name:         id.osType,
			Channel:      id.osChannel,
		}
		cfg, err := charmhub.RefreshOne(id.instanceKey, id.id, id.revision, id.channel, base)
		if err != nil {
			return nil, internalerrors.Capture(err)
		}
		cfg, err = charmhub.AddConfigMetrics(cfg, id.metrics)
		if err != nil {
			return nil, internalerrors.Capture(err)
		}
		configs[i] = cfg
	}
	config := charmhub.RefreshMany(configs...)

	ctx, cancel := context.WithTimeout(ctx, charmhub.RefreshTimeout)
	defer cancel()

	responses, err := client.RefreshWithRequestMetrics(ctx, config, metrics)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	results := make([]charmhubResult, len(responses))
	for i, response := range responses {
		if results[i], err = w.refreshResponseToCharmhubResult(response); err != nil {
			return nil, internalerrors.Capture(err)
		}
	}
	return results, nil
}

func (w *revisionUpdateWorker) getCharmhubClient(ctx context.Context) (CharmhubClient, error) {
	// Get a new downloader, this ensures that we've got a fresh
	// connection to the charm store.
	httpClient, err := w.config.NewHTTPClient(ctx, w.config.HTTPClientGetter)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	config, err := w.config.ModelConfigService.ModelConfig(ctx)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}
	charmhubURL, _ := config.CharmHubURL()

	return w.config.NewCharmhubClient(httpClient, charmhubURL, w.config.Logger)
}

// buildMetricsMetadata returns a map containing metadata key/value pairs to
// send to the charmhub for tracking metrics.
func (w *revisionUpdateWorker) buildMetricsMetadata(ctx context.Context, buildTelemetry bool) (charmhub.Metrics, error) {
	if buildTelemetry {
		return nil, nil
	}

	metrics, err := w.config.ModelService.GetModelMetrics(ctx)
	if err != nil {
		return nil, internalerrors.Capture(err)
	}

	model := metrics.Model

	return charmhub.Metrics{
		charmmetrics.Controller: {
			charmmetrics.JujuVersion: version.Current.String(),
			charmmetrics.UUID:        model.ControllerUUID.String(),
		},
		charmmetrics.Model: {
			charmmetrics.UUID:     model.UUID.String(),
			charmmetrics.Cloud:    model.Cloud,
			charmmetrics.Provider: model.CloudType,
			charmmetrics.Region:   model.CloudRegion,

			charmmetrics.NumApplications: strconv.Itoa(metrics.ApplicationCount),
			charmmetrics.NumMachines:     strconv.Itoa(metrics.MachineCount),
			charmmetrics.NumUnits:        strconv.Itoa(metrics.UnitCount),
		},
	}, nil
}

// refreshResponseToCharmhubResult converts a raw RefreshResponse from the
// charmhub API into a charmhubResult.
func (w *revisionUpdateWorker) refreshResponseToCharmhubResult(response transport.RefreshResponse) (charmhubResult, error) {
	if response.Error != nil {
		return charmhubResult{}, internalerrors.Errorf("charmhub error %s: %s", response.Error.Code, response.Error.Message)
	}

	now := w.config.Clock.Now()

	// Locate and extract the essential metadata.
	metadata, err := repository.EssentialMetadataFromResponse(response.Name, response)
	if err != nil {
		return charmhubResult{}, internalerrors.Capture(err)
	}

	var resources []resource.Resource
	for _, r := range response.Entity.Resources {
		fingerprint, err := resource.ParseFingerprint(r.Download.HashSHA384)
		if err != nil {
			w.config.Logger.Warningf("invalid resource fingerprint %q: %v", r.Download.HashSHA384, err)
			continue
		}
		typ, err := resource.ParseType(r.Type)
		if err != nil {
			w.config.Logger.Warningf("invalid resource type %q: %v", r.Type, err)
			continue
		}
		res := resource.Resource{
			Meta: resource.Meta{
				Name:        r.Name,
				Type:        typ,
				Path:        r.Filename,
				Description: r.Description,
			},
			Origin:      resource.OriginStore,
			Revision:    r.Revision,
			Fingerprint: fingerprint,
			Size:        int64(r.Download.Size),
		}
		resources = append(resources, res)
	}
	return charmhubResult{
		name:              response.Name,
		essentialMetadata: metadata,
		timestamp:         now,
		revision:          response.Entity.Revision,
		resources:         resources,
	}, nil
}

func (w *revisionUpdateWorker) scopedContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(w.catacomb.Context(context.Background()))
}

func encodeCharmURL(app application.RevisionUpdaterApplication) (*charm.URL, error) {
	// We only support charmhub charms, anything else is an error.
	if app.CharmLocator.Source != applicationcharm.CharmHubSource {
		return nil, internalerrors.Errorf("unsupported charm source %v", app.CharmLocator.Source)
	}

	arch, err := encodeArchitecture(app.CharmLocator.Architecture)
	if err != nil {
		return nil, internalerrors.Errorf("encoding architecture: %w", err)
	}

	return &charm.URL{
		Schema:       charm.CharmHub.String(),
		Name:         app.CharmLocator.Name,
		Revision:     app.CharmLocator.Revision,
		Architecture: arch,
	}, nil
}

func jitter(period time.Duration) time.Duration {
	return retry.ExpBackoff(period, period*2, 2, true)(0, 1)
}

func encodeCharmhubID(app application.RevisionUpdaterApplication, modelTag names.ModelTag) (charmhubID, error) {
	appTag, err := names.ParseApplicationTag(app.Name)
	if err != nil {
		return charmhubID{}, internalerrors.Errorf("parsing application tag: %w", err)
	}

	origin := app.Origin
	risk, err := encodeRisk(origin.Channel.Risk)
	if err != nil {
		return charmhubID{}, internalerrors.Errorf("encoding channel risk: %w", err)
	}

	channel, err := charm.MakeChannel(origin.Channel.Track, risk, origin.Channel.Branch)
	if err != nil {
		return charmhubID{}, internalerrors.Errorf("making channel: %w", err)
	}

	arch, err := encodeArchitecture(origin.Platform.Architecture)
	if err != nil {
		return charmhubID{}, internalerrors.Errorf("encoding architecture: %w", err)
	}

	osType, err := encodeOSType(origin.Platform.OSType)
	if err != nil {
		return charmhubID{}, internalerrors.Errorf("encoding os type: %w", err)
	}

	return charmhubID{
		id:          origin.ID,
		revision:    origin.Revision,
		channel:     channel.String(),
		osType:      osType,
		osChannel:   origin.Platform.Channel,
		arch:        arch,
		instanceKey: charmhub.CreateInstanceKey(appTag, modelTag),
	}, nil
}

func encodeArchitecture(a architecture.Architecture) (string, error) {
	switch a {
	case architecture.AMD64:
		return arch.AMD64, nil
	case architecture.ARM64:
		return arch.ARM64, nil
	case architecture.PPC64EL:
		return arch.PPC64EL, nil
	case architecture.S390X:
		return arch.S390X, nil
	case architecture.RISCV64:
		return arch.RISCV64, nil
	default:
		return "", internalerrors.Errorf("unsupported architecture %v", a)
	}
}

func encodeOSType(t application.OSType) (string, error) {
	switch t {
	case application.Ubuntu:
		return strings.ToLower(ostype.Ubuntu.String()), nil
	default:
		return "", internalerrors.Errorf("unsupported OS type %v", t)
	}
}

func encodeRisk(r application.ChannelRisk) (string, error) {
	switch r {
	case application.RiskStable:
		return internalcharm.Stable.String(), nil
	case application.RiskCandidate:
		return internalcharm.Candidate.String(), nil
	case application.RiskBeta:
		return internalcharm.Beta.String(), nil
	case application.RiskEdge:
		return internalcharm.Edge.String(), nil
	default:
		return "", internalerrors.Errorf("unsupported risk %v", r)
	}
}

type appInfo struct {
	id       string
	charmURL *charm.URL
}

// charmhubID holds identifying information for several charms for a
// charmhubLatestCharmInfo call.
type charmhubID struct {
	id        string
	revision  int
	channel   string
	osType    string
	osChannel string
	arch      string
	metrics   map[charmmetrics.MetricValueKey]string
	// instanceKey is a unique string associated with the application. To assist
	// with keeping KPI data in charmhub. It must be the same for every charmhub
	// Refresh action related to an application.
	instanceKey string
}

type latestCharmInfo struct {
	url       *charm.URL
	timestamp time.Time
	revision  int
	resources []resource.Resource
	appID     string
}

// charmhubResult is the type charmhubLatestCharmInfo returns: information
// about a charm revision and its resources.
type charmhubResult struct {
	name              string
	essentialMetadata corecharm.EssentialMetadata
	timestamp         time.Time
	revision          int
	resources         []resource.Resource
}
