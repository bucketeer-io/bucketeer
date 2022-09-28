import logging


class Formatter(logging.Formatter):
    def format(self, record):
        logmsg = super(Formatter, self).format(record)
        return {"msg": logmsg, "args": record.args}
