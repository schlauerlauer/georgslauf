<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');

if(!$Ende) {
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  if(isset($_POST['gruppe'])) {
    echo '<h3>von Gruppen id '.$_POST['gruppe'].'</h3>';
    echo '<ul data-role="listview" data-theme="a" data-inset="true">';
    if ($stmt = $mysqli->prepare("SELECT p.id, p.kurz, p.name, punkte.points
FROM posten p
LEFT JOIN punkte ON p.id = punkte.an
WHERE (von = ? OR von IS NULL)
ORDER BY kurz ASC")) {
      $stmt->bind_param('s', $_POST['gruppe']);
      $stmt->execute();
      $stmt->store_result();
      $stmt->bind_result($id, $kurz, $name, $punkte);
      while ($stmt->fetch()) {
        echo '<li><h2>'.$kurz.' '.$name.'</h2><p class="ui-li-aside">
        <input von="'.$_POST['gruppe'].'" an="'.$id.'" class="ppunkte" style="max-width:60px;" data-mini="true" type="number" min="0" max="100" value="'.$punkte.'"/>
        </p></li>';
      }
    }
    echo '</ul>';
  }
  else {
  echo '<h2>Postenbewertungen</h2>';
  echo '<div class="ui-field-contain"><label for="auswahl">Gruppen</label><select name="auswahl" id="auswahl" data-mini="true">';
  if ($stmt = $mysqli->prepare("SELECT id, kurz, name FROM `gruppen` ORDER BY kurz ASC")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($id, $kurz, $name);
    while ($stmt->fetch()) {
      echo '<option value="'.$id.'">'.$kurz.' - '.$name.'</option>';
    }
  }
  echo '</select></div>';
  echo '<div id="punkte">keine Gruppe ausgewählt</div>';
}
}
} else echo "Der Lauf ist beendet.";
?>
