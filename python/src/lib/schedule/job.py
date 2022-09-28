from enum import Enum


class Status(Enum):
    SUCCESS = 1
    FAIL = 2


class Job:
    def __init__(
        self,
        name,
        func,
        month="*",
        day="*",
        day_of_week="*",
        hour="*",
        minute="*",
        second="*",
    ):

        self.name = name
        self.func = func
        self.month = month
        self.day = day
        self.day_of_week = day_of_week
        self.hour = hour
        self.minute = minute
        self.second = second
