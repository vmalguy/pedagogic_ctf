import os


def init(path, randomize, file_challenge_name=None):
    with open(os.path.join(path, 'secret'), "w") as secret:
        secret.write(randomize)
