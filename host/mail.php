<?php

function sendmail($sub, $txt) {
    $to = "info@georgslauf.de";
    $head = "From: info@georgslauf.de";
    $txt = wordwrap($txt,70);
    mail($to, $sub, $txt, $head);
}
?>
