<?php
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Alle Logins</h2>";
echo '<div data-role="collapsibleset" data-theme="b" data-filter="true" data-filter-placeholder="Logins Filtern" data-content-theme="b">';
if ($stmt = $mysqli->prepare("SELECT id, username, password, rolle FROM login WHERE username != 'Janek' ORDER BY id asc")) {
  $stmt->execute();
  $stmt->store_result();
  $stmt->bind_result($id, $username, $password, $rolle);
  while ($stmt->fetch()) {
    echo '<div data-role="collapsible"><h3>'.$username.' - '.$Rollen[$rolle].'</h3><p>'.$password.'</p></div>';
  }
}
echo "</div>";

}
?>
