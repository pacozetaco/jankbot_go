JankBot is a multipurpose Discord bot.
    - Casino
        - Blackjack
        - Deathroll
        - HiLo
    -*Table Games
        -*Baccarat
    - Music Bot
    - Ark Server Monitor / Chat
    - AI chat with a 980 GTX

* = have yet to implement, in any form of bot


Notes:

-Torn between event handlers for each module or to have a switch statement in the main command handler.
    - If using event handlers for each mod, it would make each module read each event. If using a switch statement, it would make the main command handler more complex as the code grows. 
    - I think its a good idea to have a event handler for the main mods. But use a switch statement in the submods (ie. text based casino commands use a switch statement in the casino handlers.)
    - JukeBox needs a event handler.
    - Ark Server needs a handler for the chat.
    - AI chat will have a handler.
    - Command based casino will have a handler.
        - each game is in a switch statement in casino module.
    -Table games will not need a handler, all the interactions will be through view buttons. No commands neccesary. 

- Check into onMessage or onMesaageCreate... or just read the docs...