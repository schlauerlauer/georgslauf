<?php
include_once $_SERVER['DOCUMENT_ROOT'].'/includes/connect_gl.php';

$current_posten_pro_kategorie;

if ($stmt = $mysqli->prepare("SELECT sum(case when kategorie = '0' then 1 else 0 end) k0, sum(case when kategorie = '1' then 1 else 0 end) k1, sum(case when kategorie = '2' then 1 else 0 end) k2, sum(case when kategorie = '3' then 1 else 0 end) k3, sum(case when kategorie = '4' then 1 else 0 end) k4 FROM posten")) {
	$stmt->execute();
	$stmt->store_result();
	$stmt->bind_result($current_posten_pro_kategorie[0], $current_posten_pro_kategorie[1], $current_posten_pro_kategorie[2], $current_posten_pro_kategorie[3], $current_posten_pro_kategorie[4]);
	while ($stmt->fetch()) {}
}
?>
