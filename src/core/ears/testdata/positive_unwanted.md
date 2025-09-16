# Positive — Unwanted behaviour (If … then the <System> shall …)
If the payload signature is invalid, then the API shall return 401 and log the attempt.

If the file checksum fails, then the Updater shall abort and restore the last good version.

If the software detects an invalid DRAM configuration, then it shall abort the test and report the error.
