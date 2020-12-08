<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {

$pval = 0;
$totalS = 0;
$totalV = 0;
$Prevstamm = "";
if ($stmt = $mysqli->prepare("SELECT name, stamm, kategorie, anzahl, veggie, stufe FROM posten ORDER BY stamm ASC, stufe DESC, name")) {
  $stmt->execute();
  $stmt->store_result();
  $stmt->bind_result($name, $stamm, $kat, $val, $veg, $stufe);
  while ($stmt->fetch()) {
    if($stamm != $Prevstamm) {
      if($Prevstamm != "") echo "<br>Insgesamt ".$pval.' Posten';
      echo "<br><br><br><br>";
      $pval = 1;
      $Prevstamm = $stamm;
    } else $pval++;
    echo $stamm.' '.$stufe.' '.$Kat[$kat]. ' - '.$name.' -- mit '.$val.' Leitern (davon '.$veg.' Vegetarier)<br>';
  }
  echo "<br>Insgesamt ".$pval.' Posten';
}

}
?>
