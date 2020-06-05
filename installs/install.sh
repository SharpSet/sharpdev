#/bin/bash

# Uninstall any previous versions.
echo "Checking for any previous version..."
sudo rm -r /usr/local/bin/sharpdev

# Download and unpack
wget https://github.com/Sharpz7/sharpdev/releases/download/XXXXX/linux.tar.gz
sudo tar -C /usr/local/bin/ -zxvf linux.tar.gz
rm -r linux.tar.gz

# Permissions
chmod u+x /usr/local/bin/sharpdev

echo ""
echo "SHARPDEV IS NOW INSTALLED"
echo "======================="
echo "Do sharpdev for more info!"