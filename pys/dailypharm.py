import requests
from bs4 import BeautifulSoup

url = 'http://www.dailypharm.com/Users/News/NewsList.html'
req = requests.get(url)
bs = BeautifulSoup(req.content, 'lxml')

wrapper = bs.find('div', class_='seachResult').find('ul', recursive=False)
items = wrapper.find_all('li', class_='newsList')

for item in items[:10]:
    a_tag = item.find('a')
    href = 'http://www.dailypharm.com/Users/News/' + a_tag['href']
    title = item.find('div', class_='listHead').find(text=True).strip()
    date = item.find('div', class_='listHead').find('span', class_='newsDate').get_text().strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.text, 'lxml')

    wrapper2 = bs2.find('div', class_='newsContents')
    contents = ''
    for it in wrapper2.find_all(text=True, recursive=False):
        try:
            contents += it.strip().encode('iso-8859-1').decode('cp949')
        except:
            continue

    print('==========================')
    print(title)
    print(href)
    print(date)
    print(contents)