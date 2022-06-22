<?php
include_once '/var/www/vhosts/hosting101172.af98b.netcup.net/www/georgslauf/dev/includes/connect_gl.php';
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Whatsapp an Posten schreiben</h2>";

  $val = 0;
  echo '<div class="ui-grid-d ui-responsive">';
  if ($stmt = $mysqli->prepare("SELECT kurz, kontakt FROM posten order by kurz")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($eins, $zwei);
    while ($stmt->fetch()) {
      $val++;
      $bs;
      switch ($val) {
        case 1: $bs = 'a';
          break;
        case 2: $bs = 'b';
          break;
        case 3: $bs = 'c';
          break;
        case 4: $bs = 'd';
          break;
        case 5: $bs = 'e';
          $val = 0;
          break;
      }
      echo '<div class="ui-block-'.$bs.'"><a class="ui-btn ui-icon-comment ui-btn-icon-left" href="https://api.whatsapp.com/send?phone=49'.substr($zwei, 1).'" type="button">'.$eins.'</a></div>';
    }
  }
  echo '</div>';
} else echo "Keine Berechtigung";
?>
