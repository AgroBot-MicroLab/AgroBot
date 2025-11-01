# Setting up and running
1. First of all make sure you have docker installed by opening the console and doing docker --version.
2. Go to the project root and run `docker compose up`, this should start the whole project
3. After the project started create a hotspot on your phone and set up the Raspberry PI in headless mode so that it connects to your phone's hotspot. This can be done by downloading the Raspberry PI flasher software and flashing an SD card.
**IMPORTANT:** make sure the wifi you are using is 2.4GHz to save your sanity and time.
4. After that you can ssh to the Raspberry PI and move the drone/ diretory to it, using scp.
5. After you moved it to the PI, run `pip install -r requirements.txt`. You might also need to do some `sudo apt install`s, to enable the camera. 
6. After everything is set up run `python main.py sub --host broker.emqx.io --port 1883 --topic "agro/2f37308d-ccf8-4685-8329-c164599259ca/cmd" --qos 1`
7. Go back to your system an install and setup ardupilot by cloning this repo: https://github.com/ArduPilot/ardupilot
8. You would have to do something like `./waf configure --board sitl` then `./waf rover` and then you would be able to run it with
`./Tools/autotest/sim_vehicle.py -v Rover --console --map --out=udp:127.0.0.1:14550 --custom-location=47.061657183060966,28.867524495508608,10,0`

After that the project should be set up and running.
