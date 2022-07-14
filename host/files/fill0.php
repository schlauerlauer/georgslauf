<?php
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 3) {

  if ($stmt = $mysqli->prepare("INSERT IGNORE INTO punkte (von, an, points)
SELECT DISTINCT posten.id, gruppen.id, 0
FROM gruppen, posten")) {
    $stmt->execute();
    echo "Alle Bewertungen der Posten mit Null gefüllt.";
  }
  else echo "Beim Einfügen der Posten Punkte ist etwas schiefgelaufen!";

  if ($stmt = $mysqli->prepare("INSERT IGNORE INTO punkte (von, an, points)
SELECT DISTINCT gruppen.id, posten.id, 0
FROM gruppen, posten")) {
    $stmt->execute();
    echo "Alle Bewertungen der Gruppen mit Null gefüllt.";
  }
  else echo "Beim Einfügen der Gruppen Punkte ist etwas schiefgelaufen";

} else echo "Keine Berechtigung";
?>
