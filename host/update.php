<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/master/includes/connect_gl.php';
include_once 'settings.php';
require('../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {

if(isset($_POST['k'])) {
  if ($stmt = $mysqli->prepare("UPDATE posten SET kurz = ? WHERE id = ?")) {
    $stmt->bind_param('ss', $_POST['k'], $_POST['id']);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
}

if(isset($_POST['color'])) {
  if ($stmt = $mysqli->prepare("UPDATE posten SET color = ? WHERE id = ?")) {
    $stmt->bind_param('ss', $_POST['color'], $_POST['id']);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
}

if(isset($_POST['select'])) {
  $size = 100 / sizeof($Farben);
  echo '<table width="100%"><tr>';
  for($i = 0; $i < sizeof($Farben); $i++) {
    echo '<td width="'.$size.'%" id="'.$i.'" posten="'.$_POST['id'].'" style="background-color:'.$Farben[$i].';" class="select">&nbsp;</td>';
  }
  echo '</tr></table>';
}

if(isset($_POST['kurz'])) {
  if ($stmt = $mysqli->prepare("UPDATE gruppen SET kurz = ? WHERE id = ?")) {
    $stmt->bind_param('ss', $_POST['kurz'], $_POST['id']);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
}

if(isset($_POST['coor'])) {
  $cor = "";
  if($_POST['xy'] == "x") $cor = "x_axis";
  else $cor = "y_axis";
  if ($stmt = $mysqli->prepare("UPDATE posten SET ".$cor." = ? WHERE kurz = ?")) {
    $stmt->bind_param('ss', $_POST['coor'], $_POST['posten']);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
}

if(isset($_POST['startG'])) {
  if ($stmt = $mysqli->prepare("UPDATE posten SET startGruppen = ? WHERE kurz = ?")) {
    $stmt->bind_param('ss', $_POST['startG'], $_POST['kurzel']);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
}

}
?>
