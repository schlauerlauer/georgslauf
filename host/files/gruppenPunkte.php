<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');

if(!$Ende) {
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  if(isset($_POST['gruppe'])) {
    echo '<h3>von Posten id '.$_POST['gruppe'].'</h3>';
    echo '<ul data-role="listview" data-theme="a" data-inset="true">';
    if ($stmt = $mysqli->prepare("SELECT g.id, g.kurz, g.name, g.stufe, g.stamm, punkte.points
FROM gruppen g
LEFT JOIN punkte ON g.id = punkte.an
WHERE (von = ? OR von IS NULL)
ORDER BY kurz ASC")) {
      $stmt->bind_param('s', $_POST['gruppe']);
      $stmt->execute();
      $stmt->store_result();
      $stmt->bind_result($id, $kurz, $name, $stufe, $stamm, $punkte);
      while ($stmt->fetch()) {
        echo '<li><h2>'.$kurz.' '.$name.'</h2><p class="ui-li-aside">
        <input id="g'.$id.'" von="'.$_POST['gruppe'].'" an="'.$id.'" class="ppunkte" style="max-width:60px;" data-mini="true" type="number" min="0" max="100" value="'.$punkte.'"/>
        </p></li>';
      }
    }
    echo '</ul>';
  }
  else {
  echo '<h2>Gruppenbewertungen</h2>';
  echo '<div class="ui-field-contain"><label for="auswahl">Posten</label><select name="auswahl" id="auswahl" type="gruppe" data-mini="true">';
  if ($stmt = $mysqli->prepare("SELECT id, kurz, name, stufe, stamm FROM `posten` ORDER BY kurz ASC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($id, $kurz, $name, $stufe, $stamm);
    while ($stmt->fetch()) {
      echo '<option value="'.$id.'">'.$kurz.' - '.$name.' ('.$stufe.' '.$stamm.')</option>';
    }
  }
  echo '</select></div>';
  echo '<div id="punkte">keinen Posten ausgewählt</div>';
}
}
} else echo "Der Lauf ist beendet.";
?>
