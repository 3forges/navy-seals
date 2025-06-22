Ok, so it works great, and here is the process (very simple actually) to re-implement in telegram:

* First I need the Telegram Bot token: that, is picked up by navy seal from its configuration, it is a secret that the user needs not to know about, and that does not change for any user.
* NavySeal generates a random complex "password": something that is going to be unique per user
* The user is asked to add the telegram bot just by giving him the link:  https://telegram.me/ToutatisRoBot
* The bot polls the Telegram API until the user has sent a message consisting of only the "password":
   * there is a timeout, if timedout, thepassword expires and the user can request a new "password"
   * if the u ser actually sent the "password", then navyseal can retrieve the CHat ID (and not the user ID) : that "ChatID" is the unique ID of a conversation happening only between the desied user ad the telegram bot, so in that chat I can send a message with the QR code

Et voilà l'ami