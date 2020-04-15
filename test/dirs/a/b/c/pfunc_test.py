
def do_print():
    print "hello world"


def add(a, b):
    return a + b


def names_of_three_people(a, b, c):
    return a['name'] + " and " + b['name'] + " and " + c['name']


def divide(a, b):
    return a / b


def float_divide(a, b):
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


def first_param_and_other_params(first, **other):
    total = other
    total['first'] = first
    return total
