#!/bin/bash
# As a member of juju-qa,  Visit each the jenkins master and slaves
# and update their branches.
# passing 'true' as an arg will driect the script to try to update cloud-city.
set -eux

MASTER="juju-ci.vapour.ws"
KEY="staging-juju-rsa"
export JUJU_ENV="juju-ci3"


update_jenkins() {
    host=$1
    echo "updating $host"
    if [[ "$CLOUD_CITY" == "true" ]]; then
        bzr branch lp:~juju-qa/+junk/cloud-city \
            bzr+ssh://jenkins@$host/var/lib/jenkins/cloud-city.new
    fi
    ssh jenkins@$host << EOT
#!/bin/bash
set -eux
if [[ "$CLOUD_CITY" == "true" ]]; then
    (cd ~/cloud-city; bzr revert; cd -)
    bzr pull -d ~/cloud-city ~/cloud-city.new
    rm -r ~/cloud-city.new
fi
cd ~/juju-release-tools
bzr pull
cd ~/repository
bzr pull
cd ~/juju-ci-tools
bzr pull
make install-deps
if [[ -d ~/ci-director ]]; then
    cd ~/ci-director
    bzr pull
fi
EOT
}


CLOUD_CITY="false"
while [[ "${1-}" != "" ]]; do
    case $1 in
        --cloud-city)
            CLOUD_CITY="true"
            ;;
    esac
    shift
done

SLAVES=$(juju status '*-slave*' | grep public-address | sed -r 's,^.*: ,,')
if [[ -z $SLAVES ]]; then
    echo "Set JUJU_HOME to juju-qa's environments and switch to juju-ci."
    exit 1
fi
if [[ ! $SLAVES =~ ^.*10\.125\.0\.10.*$ ]]; then
    echo "The kvm-slave lost its machine and unit agents."
    SLAVES="$SLAVES 10.125.0.10 15.125.114.8 osx-slave.vapour.ws"
fi

SKIPPED=""
for host in $MASTER $SLAVES; do
    update_jenkins $host || SKIPPED="$SKIPPED $host"
done

if [[ -n "$SKIPPED" ]]; then
    set +x
    echo
    echo "These hosts were skipped because there was an error"
    echo "$SKIPPED"
fi

