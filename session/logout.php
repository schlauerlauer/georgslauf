<?php
session_start();
if(session_destroy()) //Destroy all sessions
{
header("Location: /");
}
?>