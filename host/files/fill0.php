<?php
include_once $__SERVER["DOCUMENT_ROOT"].'/includes/connect_gl.php';
require('../../session/session.php');

if(isset($login_session) && $_SESSION['rolle'] >= 4) {

  if ($stmt = $mysqli->prepare("INSERT IGNORE INTO punkte (von, an, points)
SELECT DISTINCT posten.id, gruppen.id, 0
FROM gruppen, posten
WHERE (posten.stufe = 'RoPo' AND gruppen.stufe >= 2)
OR (posten.stufe != 'RoPo' AND gruppen.stufe < 2)")) {
    $stmt->execute();
    echo "Alle Bewertungen der Posten mit Null gefüllt.<br />";
  }
  else echo "Beim Einfügen der Posten Punkte ist etwas schiefgelaufen!<br />";

  if ($stmt = $mysqli->prepare("INSERT IGNORE INTO punkte (von, an, points)
SELECT DISTINCT gruppen.id, posten.id, 0
FROM gruppen, posten
WHERE (posten.stufe = 'RoPo' AND gruppen.stufe >= 2)
OR (posten.stufe != 'RoPo' AND gruppen.stufe < 2)")) {
    $stmt->execute();
    echo "Alle Bewertungen der Gruppen mit Null gefüllt.<br />";
  }
  else echo "Beim Einfügen der Gruppen Punkte ist etwas schiefgelaufen<br />";

} else echo "Keine Berechtigung";
?>
