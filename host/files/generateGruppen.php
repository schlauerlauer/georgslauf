<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Alle Gruppen</h2>";
  echo "<h4>für Rückmeldung an die Stämme</h4>";
$pval = 0;
$Prevstamm = "";
if ($stmt = $mysqli->prepare("SELECT name, stufe, stamm, size, veggies FROM gruppen ORDER BY stamm, stufe, name asc")) {
  $stmt->execute();
  $stmt->store_result();
  $stmt->bind_result($name, $stufe, $stamm, $size, $veggies);
  while ($stmt->fetch()) {
    if($stamm != $Prevstamm) {
      if($Prevstamm != "") echo "<br>Insgesamt ".$pval.' Gruppen';
      echo "<br><br><br><br>";
      $pval = 1;
      $Prevstamm = $stamm;
    } else $pval++;
    echo $stamm.' '.$Stufe[$stufe].' - '.$name.' -- mit '.$size.' Kindern (davon '.$veggies.' Vegetarier)<br>';
  }
  echo "<br>Insgesamt ".$pval.' Gruppen';
}
}
?>
