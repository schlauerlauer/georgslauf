<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Startgruppen</h2>";
  if ($stmt = $mysqli->prepare("SELECT id, kurz, startGruppen FROM posten ORDER BY kurz")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($id, $kurz, $start);
    while ($stmt->fetch()) {
      echo '<label for="startGruppen">'.$kurz.'</label>
            <input type="text" class="startG" kurz="'.$kurz.'" value="'.$start.'" placeholder="Start Gruppen">';
    }
  } else echo "Etwas ist schiefgelaufen";
} else echo "Keine Berechtigung";
?>
