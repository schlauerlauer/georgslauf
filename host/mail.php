<?php

function sendmail($sub, $txt) {
    $to = "georgslauf@pfadi-fc.de";
    $head = "From: georgslauf@pfadi-fc.de";
    $txt = wordwrap($txt,70);
    mail($to, $sub, $txt, $head);
}
?>
