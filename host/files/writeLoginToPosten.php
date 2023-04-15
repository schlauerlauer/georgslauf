<?php
include_once '../settings.php';
require('../../session/session.php');
if(isset($login_session) && $_SESSION['rolle'] >= 3) {
  echo "<h2>Zugang an Posten schicken</h2>";

  $val = 0;
  echo '<div class="ui-grid-d ui-responsive">';
  if ($stmt = $mysqli->prepare("SELECT username, password, kontakt FROM login, posten WHERE posten.kurz = login.username order by kurz")) {
    $stmt->execute();
    $stmt->store_result();
    $stmt->bind_result($eins, $zwei, $drei);
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
      echo '<div class="ui-block-'.$bs.'"><a class="ui-btn ui-icon-comment ui-btn-icon-left" target="_blank" href="https://api.whatsapp.com/send?phone=49'.substr($drei, 1).'&text=Hallo%20Posten%20'.$eins.'%0D%0A%0D%0AHier%20Eure%20Zugangsdaten%20f%C3%BCr%20https%3A%2F%2Fgeorgslauf.de%0D%0A%0D%0AName%20'.$eins.'%0D%0APasswort%20'.urlencode($zwei).'%0D%0A%0D%0ABitte%20gib%20die%20Infos%20an%20die%20anderen%20Leiter%20am%20Posten%20weiter%21%0D%0A%0D%0AViele%20Gr%C3%BC%C3%9Fe%2C%0D%0AAnsgar%20%3A%29" type="button">'.$eins.'</a></div>';
      //Nur Passwort echo '<div class="ui-block-'.$bs.'"><a class="ui-btn ui-icon-comment ui-btn-icon-left" href="https://api.whatsapp.com/send?phone=49'.substr($drei, 1).'&text='.$zwei.'" type="button">'.$eins.'</a></div>';
    }
  }
  echo '</div>';
} else echo "Keine Berechtigung";
?>
