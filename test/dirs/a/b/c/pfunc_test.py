
def func1():
    print "this is func1"


def func2(a, b):
    return a + b


def func3(a, b, c):
    return a['name'] + " and " + b['name'] + " and " + c['name']


def func4(a, b):
    return a / b


def func5(a, b):
    return float(a) / float(b)


def func_return_struct(name, age, hobby1, hobby2):
    return {
        "name": name,
        "age": age,
        "hobby": [
            hobby1,
            hobby2
        ]
    }