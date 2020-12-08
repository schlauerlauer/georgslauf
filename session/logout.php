<?php
if(session_destroy()) //Destroy all sessions
{
header("Location: /");
}
?>