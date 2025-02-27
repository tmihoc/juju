run_hook_dispatching_script() {
	echo

	model_name="test-hook-dispatching"
	file="${TEST_DIR}/${model_name}.log"

	ensure "${model_name}" "${file}"

	# the log messages the test looks for do not appear if root
	# log level is WARNING.
	juju model-config logging-config="<root>=INFO"

	# shellcheck disable=SC2046
	juju deploy $(pack_charm ./testcharms/charms/ubuntu-plus)
	wait_for "ubuntu-plus" "$(idle_condition "ubuntu-plus")"

	juju debug-log --include unit-ubuntu-plus-0 | grep -q "via hook dispatching script: dispatch" || true

	juju run ubuntu-plus/0 no-dispatch filename=test-dispatch

	# wait for update-status
	wait_for "Hello from update-status" "$(workload_status ubuntu-plus 0).message"

	# check it was not via dispatch
	juju debug-log --include unit-ubuntu-plus-0 | grep -q 'ran "update-status" hook (via explicit, bespoke hook script)' || true

	destroy_model "${model_name}"
}

test_dispatching_script() {
	if [ "$(skip 'test_dispatching_script')" ]; then
		echo "==> TEST SKIPPED: dispatch"
		return
	fi

	(
		set_verbosity

		cd .. || exit

		run "run_hook_dispatching_script"
	)
}
