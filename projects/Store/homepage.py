# This Python file uses the following encoding: utf-8
import os
from pathlib import Path
import sys
import requests
import time
import json
from PySide6.QtWidgets import QApplication, QMainWindow
from PySide6.QtCore import QFile, Qt
from PySide6.QtUiTools import QUiLoader
from PySide6.QtCore import Slot
from PySide6.QtGui import QMovie

try:
    os.chdir(sys._MEIPASS)
    print(sys._MEIPASS)
except:
    pass

class HomePage(QMainWindow):
    def __init__(self):
        super(HomePage, self).__init__()
        self.load_ui()
        self.baseurl = "http://localhost:8000/"
#        self.baseurl = "http://192.168.64.2:30811/"

        self.isConnected = False
        self.check_connection()
        self.w.networkButton.clicked.connect(self.check_connection)

        self.w.pageTwoEnter.clicked.connect(self.search_item)
        self.w.pageThreeEnter.clicked.connect(self.create_item)
        self.w.pageFourEnter.clicked.connect(self.delete_item)

    @Slot()
    def search_item(self):
        if self.isConnected == True:
            self.w.pageTwoResults.setText("")
            url = self.baseurl + "search?title=" + self.w.itemSearch.text()
            r = requests.request(method='GET',url=url)
            games = r.json()
            results = ""

            for game in games:
                if game["error"] is not None:
                    results += "Error: " + game["error"] + '\n'
                    if game["code"] is not None:
                        results += "Error Code: " + game["code"] + '\n'

                    results+= '\n'
                    continue

                results += ("Name: " + game["name"] + '\n' +
                            "Developer: " + ",".join(game["companies"]) + '\n' +
                            "Release Date: " + game["date"] + '\n' +
                            "Critic Ratings: " + game["ratings"] + '%\n' +
                            "Platforms: " + ",".join(game["platforms"]) + '\n' +
                            "Price: " + str(game["price"]) + '\n' +
                            "Quantity: " + str(game["quantity"]) + '\n' +
                            "In Stock: " + str(game["in_stock"]) + '\n' +
                            "SKU Number: " + str(game["sku"]) + '\n\n' )



            self.w.pageTwoResults.setText(results)
            self.w.pageTwoResults.adjustSize()
            self.w.pageTwoResults.show()

            self.w.itemSearch.setText("Enter a title")
        else:
            self.check_connection()

    @Slot()
    def create_item(self):
        if self.isConnected == True:
            url = (self.baseurl + "create?title=" + self.w.itemAdd.text().replace(" ","_") + "&price=" + self.w.priceAdd.text() + "&quantity=" + self.w.quantityAdd.text())
            r = requests.request(method='POST',url=url)

            self.w.itemAdd.setText("Enter a title")
            self.w.priceAdd.setText("Enter a price")
            self.w.quantityAdd.setText("Enter a quantity")
        else:
            self.check_connection()

    @Slot()
    def delete_item(self):
        if self.isConnected == True:
            url = (self.baseurl + "delete?title=" + self.w.itemDelete.text().replace(" ","_"))
            r = requests.request(method='POST',url=url)

            self.w.itemDelete.setText("Enter a title")
        else:
            self.check_connection()

    @Slot()
    def check_connection(self):
        try:
            requests.get(self.baseurl).status_code
            self.w.networkLabel.move(485,315)
            self.w.networkButton.hide()
            self.w.networkLabel.setText("Network Connection Established!")
            self.w.pageTwoEnter.setText("Enter")
            self.w.pageThreeEnter.setText("Enter")
            self.w.pageFourEnter.setText("Enter")
            self.isConnected = True
        except:
            self.w.networkLabel.move(415,315)
            self.w.networkLabel.setText("Network Connection Not Established! Try Again.")
            self.w.pageTwoEnter.setText("Network Not Connected!")
            self.w.pageThreeEnter.setText("Network Not Connected!")
            self.w.pageFourEnter.setText("Network Not Connected!")
            self.w.networkButton.show()

    def init_pages(self):
       self.w.networkButton.hide()
       self.w.pageTwoResults.hide()

       scriptDir = os.path.dirname(os.path.realpath(__file__))
#       self.gif_path = os.path.join(scriptDir,'images','gaming.gif')
       self.gif_path = 'gaming.gif'
       self.gif = QMovie(self.gif_path)
       self.w.gameGif.setMovie(self.gif)
       self.gif.start()

    def load_ui(self):
        loader = QUiLoader()
        path = os.fspath(Path(__file__).resolve().parent / "homepage.ui")
        ui_file = QFile(path)
        ui_file.open(QFile.ReadOnly)
        self.w = loader.load(ui_file, self)
        self.init_pages()
        self.w.show()
        ui_file.close()

if __name__ == "__main__":
    app = QApplication([])
    home_page = HomePage()
    sys.exit(app.exec())
