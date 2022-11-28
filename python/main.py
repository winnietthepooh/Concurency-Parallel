import asyncio
from time import perf_counter

import aiohttp

BASE_URL = "https://api.hypixel.net/skyblock/auctions"


async def fetch(session: aiohttp.ClientSession, page: int):
    async with session.get(f'{BASE_URL}?page={page}') as req:
        if req.status != 200:
            req.raise_for_status()
        return await req.json()


async def fetch_all(session: aiohttp.ClientSession, pages: list[int]):
    tasks = []
    for page in pages:
        tasks.append(asyncio.create_task(fetch(session, page)))
    return await asyncio.gather(*tasks)


async def main():
    async with aiohttp.ClientSession() as session:
        init = await fetch(session, 0)
        if not init["success"]:
            return
        pages = []
        pages.extend(range(1, init["totalPages"], 1))
        await fetch_all(session, pages)


if __name__ == "__main__":
    start = perf_counter()
    asyncio.run(main())
    stop = perf_counter()
    print("time taken:", stop - start)
