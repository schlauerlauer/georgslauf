<?php
include_once '../../../includes/connect_gl.php';

if(isset($_POST['p'])) {
  $user = 'wX';
  $pVal = $_POST['p'];
  if ($pVal > 100) $pVal = 100;
  else if ($pVal < 0) $pVal = 0;
  $vergleich = "< 2";
  if(substr($user,0,1) == "r") $vergleich = ">= 2";
  $pID = 114; //TODO beim login abfragen und setzen!
  if ($stmt = $mysqli->prepare("INSERT INTO punkte (von, an, points) VALUES (?, IF((SELECT stufe FROM gruppen WHERE id = ?)".$vergleich.", ?, -1), ?)
  ON DUPLICATE KEY UPDATE points=VALUES(points), time=CURRENT_TIMESTAMP")) {
    $stmt->bind_param('iiii', $pID, $_POST['id'], $_POST['id'], $pVal);
    $stmt->execute();
    echo "ok";
  }
  else echo "error";
} else echo "error";
 ?>
