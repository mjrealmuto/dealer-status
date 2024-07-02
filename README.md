# dealer-status
This Go Application is used to query sites and determine whether or not they are hosted by DealerInspire ( Cars Commerce ).
The application will create a SSH Tunnel into DealerInspire which will then allow to query the Dealer Management table 
for active accounts.  We then take that information and, concurretnly, run CNAME checks on the site URLs.  If it is 
determined that the site is not hosted by DealerInspire ( Cars Commerce ) we will write those sites to a CSV in
the `assets` folder of this project.

## Set-up
There will be a few files we'll have to edit before we can get this fully working. 
- The first thing we'll need to do is take the `.env.example` file and duplicate that. Rename it to `.env`. Inside 
the file there will be some values that will need to be populated.  These fields are DB and SSH information.  
The DB fields will be what we use to connect to the Production Dealer Management Database.  If you need 
some idea as to what these are, please feel free to contact me.  The SSH prefixed one are what we use
to tunnel in when we connect to the database.  
- Lastly, we'll edit the `Makefile` and put in the path to your private RSA key. Why do we need to do this? 
When we connect to the database, we tunnel in via SSH so we are going to mount our `id_rsa` file inside
the Container so we can connect.

Once this is all complete, you will be able to just run `make` in the root of the application.  Once the 
application is up and running you'll get some information in regards to the status of the SSH Tunnel
and DB Ping.  Once that is done, the process to get the report will run between 5-10 minutes.  The
Result will be put in the `assets` folder of this project, as stated earlier.