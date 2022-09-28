import pymysql.cursors


class Client:
    def __init__(
        self,
        mysql_user: str,
        mysql_pass: str,
        mysql_host: str,
        mysql_port: int,
        mysql_db_name: str,
    ):
        self.mysql_user = mysql_user
        self.mysql_pass = mysql_pass
        self.mysql_host = mysql_host
        self.mysql_port = mysql_port
        self.mysql_db_name = mysql_db_name

    def get_conn(self):
        return pymysql.connect(
            host=self.mysql_host,
            user=self.mysql_user,
            password=self.mysql_pass,
            database=self.mysql_db_name,
            cursorclass=pymysql.cursors.DictCursor,
        )
