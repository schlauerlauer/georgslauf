<?php
include_once 'settings.php';

function sendmail($sub, $txt) {
    $to = $Email;
    $head = "From: $Email";
    $txt = wordwrap($txt,70);
    mail($to, $sub, $txt, $head);
}
?>
