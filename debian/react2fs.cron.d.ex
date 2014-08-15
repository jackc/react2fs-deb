#
# Regular cron jobs for the react2fs package
#
0 4	* * *	root	[ -x /usr/bin/react2fs_maintenance ] && /usr/bin/react2fs_maintenance
