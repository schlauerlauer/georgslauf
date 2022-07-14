<?php
include_once '../settings.php';
require('../../session/session.php');

if($Backup) {
  $filename = 'db_backup_'.date(ymdaG).'.sql';
  echo "In Arbeit";
} else echo "Die Backup Funktion ist ausgeschaltet";
?>
