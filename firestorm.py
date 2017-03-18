#!/usr/bin/env python

from pprint import pprint
import aiohttp
import random
import asyncio
import time

class Result():
    def __init__(self, method, path, status, time):
        self.method = method
        self.path = path
        self.status = status
        self.time = time

    def __str__(self):
        return '[{}] {} "{}" ({}ms)'.format(
            self.status, self.method, self.path, self.time 
        )

class FirestormClient():
    def __init__(self, tasks, min_delay=1000, max_delay=1000, options=None):
        self.tasks = tasks
        self.min_delay = min_delay
        self.max_delay = max_delay
        self.options = options
        self.active = True

    async def exec(self, firestorm, host):
        self.session = aiohttp.ClientSession(loop=firestorm.loop, **self.options)
        self.firestorm = firestorm
        self.host = host
        while self.active:
            task = random.choice(self.tasks)
            html = task(self)
            sleep_time_ms = random.randint(self.min_delay, self.max_delay)
            await asyncio.sleep(sleep_time_ms / 1000)

    async def fetcher(self, method, path, **kwargs):
        start = time.time()
        response = await self.session.request(method, self.host + path, **kwargs)
        response.close()
        end = time.time()
        self.firestorm.log('GET', path, response, end - start)
        return response

    def fetch(self, *args, **kwargs):
        asyncio.ensure_future(self.fetcher(*args, **kwargs))

class Firestorm():
    def __init__(self, host, options=None, clients=10, tasks=None):
        self.host = host
        self.history = []
        self.loop = asyncio.get_event_loop()
        self.clients = [FirestormClient(tasks, options=options) for i in range(clients)]

    def fire(self):
        self.loop.run_until_complete(self.run())

    def log(self, method, path, response, time_in_seconds):
        record = Result(
            method=method,
            path=path,
            status=response.status,
            time=time_in_seconds * 1000,
        )
        self.history.append(record)
        print(record)

    async def run(self):
        futures = []
        for client in self.clients:
            futures.append(client.exec(self, self.host))
        try:
            await asyncio.wait(futures)
        except KeyboardInterrupt:
            for client in self.clients:
                client.active = False
        await asyncio.wait(futures)
        print('Summary:')
        print('{} requests sent'.format(
            len(self.history)
        ))


if __name__ == '__main__':
    def example_task(client):
        print('Running example task')
        client.fetch('GET', '/')
    Firestorm('http://localhost:80', clients=1, tasks=[example_task]).fire()

