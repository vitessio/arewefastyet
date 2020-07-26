import mysql.connector

def connectdb(hostname,username,password,database):
    mydb = mysql.connector.connect(
      host=hostname,
      user=username,
      password=password,
      database=database
    )
    return mydb

# Testing 
#connectdb("localhost","akilan","akilan","vitess_benchmark")
