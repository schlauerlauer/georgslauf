<?php

function sendmail($sub, $txt) {
    $to = "gl20@stamm-prm.de";
    $head = "From: gl20@stamm-prm.de";
    $txt = wordwrap($txt,70);
    mail($to, $sub, $txt, $head);
}
?>
