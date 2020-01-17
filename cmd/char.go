package cmd

func getLogo() string {
	Shuffle(logos)
	return logos[0]
}

var logos = []string{
	`
	
          _____                _____                    _____          
         /\    \              /\    \                  /\    \         
        /::\    \            /::\    \                /::\    \        
       /::::\    \           \:::\    \              /::::\    \       
      /::::::\    \           \:::\    \            /::::::\    \      
     /:::/\:::\    \           \:::\    \          /:::/\:::\    \     
    /:::/__\:::\    \           \:::\    \        /:::/  \:::\    \    
   /::::\   \:::\    \          /::::\    \      /:::/    \:::\    \   
  /::::::\   \:::\    \        /::::::\    \    /:::/    / \:::\    \  
 /:::/\:::\   \:::\____\      /:::/\:::\    \  /:::/    /   \:::\    \ 
/:::/  \:::\   \:::|    |    /:::/  \:::\____\/:::/____/     \:::\____\
\::/   |::::\  /:::|____|   /:::/    \::/    /\:::\    \      \::/    /
 \/____|:::::\/:::/    /   /:::/    / \/____/  \:::\    \      \/____/ 
       |:::::::::/    /   /:::/    /            \:::\    \             
       |::|\::::/    /   /:::/    /              \:::\    \            
       |::| \::/____/    \::/    /                \:::\    \           
       |::|  ~|           \/____/                  \:::\    \          
       |::|   |                                     \:::\    \         
       \::|   |                                      \:::\____\        
        \:|   |                                       \::/    /        
         \|___|                                        \/____/         
                                                                       

	`,
	`

         _         _             _      
        /\ \      /\ \         /\ \     
       /  \ \     \_\ \       /  \ \    
      / /\ \ \    /\__ \     / /\ \ \   
     / / /\ \_\  / /_ \ \   / / /\ \ \  
    / / /_/ / / / / /\ \ \ / / /  \ \_\ 
   / / /__\/ / / / /  \/_// / /    \/_/ 
  / / /_____/ / / /      / / /          
 / / /\ \ \  / / /      / / /________   
/ / /  \ \ \/_/ /      / / /_________\  
\/_/    \_\/\_\/       \/____________/  
                                        
               

	`,
	`

      ___           ___           ___     
     /\  \         /\  \         /\  \    
    /::\  \        \:\  \       /::\  \   
   /:/\:\  \        \:\  \     /:/\:\  \  
  /::\~\:\  \       /::\  \   /:/  \:\  \ 
 /:/\:\ \:\__\     /:/\:\__\ /:/__/ \:\__\
 \/_|::\/:/  /    /:/  \/__/ \:\  \  \/__/
    |:|::/  /    /:/  /       \:\  \      
    |:|\/__/     \/__/         \:\  \     
    |:|  |                      \:\__\    
     \|__|                       \/__/    

	`,
	`

      ___                         ___     
     /\  \                       /\__\    
    /::\  \         ___         /:/  /    
   /:/\:\__\       /\__\       /:/  /     
  /:/ /:/  /      /:/  /      /:/  /  ___ 
 /:/_/:/__/___   /:/__/      /:/__/  /\__\
 \:\/:::::/  /  /::\  \      \:\  \ /:/  /
  \::/~~/~~~~  /:/\:\  \      \:\  /:/  / 
   \:\~~\      \/__\:\  \      \:\/:/  /  
    \:\__\          \:\__\      \::/  /   
     \/__/           \/__/       \/__/    

	`,
	`

      ___                       ___     
     /  /\          ___        /  /\    
    /  /::\        /  /\      /  /:/    
   /  /:/\:\      /  /:/     /  /:/     
  /  /:/~/:/     /  /:/     /  /:/  ___ 
 /__/:/ /:/___  /  /::\    /__/:/  /  /\
 \  \:\/:::::/ /__/:/\:\   \  \:\ /  /:/
  \  \::/~~~~  \__\/  \:\   \  \:\  /:/ 
   \  \:\           \  \:\   \  \:\/:/  
    \  \:\           \__\/    \  \::/   
     \__\/                     \__\/    

	`,
	`

      ___                         ___     
     /  /\          ___          /  /\    
    /  /::\        /__/\        /  /::\   
   /  /:/\:\       \  \:\      /  /:/\:\  
  /  /::\ \:\       \__\:\    /  /:/  \:\ 
 /__/:/\:\_\:\      /  /::\  /__/:/ \  \:\
 \__\/~|::\/:/     /  /:/\:\ \  \:\  \__\/
    |  |:|::/     /  /:/__\/  \  \:\      
    |  |:|\/     /__/:/        \  \:\     
    |__|:|~      \__\/          \  \:\    
     \__\|                       \__\/    

	`,
	`

___________________________________________        
 ___________________________________________       
  ___________________/\\\____________________      
   __/\\/\\\\\\\___/\\\\\\\\\\\_____/\\\\\\\\_     
    _\/\\\/////\\\_\////\\\////____/\\\//////__    
     _\/\\\___\///_____\/\\\_______/\\\_________   
      _\/\\\____________\/\\\_/\\__\//\\\________  
       _\/\\\____________\//\\\\\____\///\\\\\\\\_ 
        _\///______________\/////_______\////////__

	`,
}