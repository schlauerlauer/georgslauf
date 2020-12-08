<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {

echo '<div data-role="collapsibleset" data-theme="b" data-filter="true" data-filter-placeholder="Logins Filtern" data-content-theme="b">';
if ($stmt = $mysqli->prepare("SELECT id, username, password, rolle FROM login WHERE rolle < 4 ORDER BY id asc")) {
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
