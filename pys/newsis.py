import requests
from bs4 import BeautifulSoup

url = 'http://www.newsis.com/realnews/'

req = requests.get(url)
bs = BeautifulSoup(req.text, 'lxml')

wrapper = bs.find('div', class_='lst_p6 mgt21')
items = wrapper.find_all('li', class_='p1_bundle')

for item in items[:1]:
    a_tag = item.find('strong', class_='title').find('a')
    title = a_tag.get_text().strip()
    href = 'http://www.newsis.com' + a_tag['href']
    date = item.find('span', class_='date').get_text().strip(' | ')[1].strip()

    req2 = requests.get(href)
    bs2 = BeautifulSoup(req2.content, 'lxml')
    wrapper2 = bs2.find('div', id='textBody')

    contents = wrapper2.get_text().strip()
    remove = wrapper2.find('div', class_='summary_view')
    remove2 = wrapper2.find_all('div', class_='view_text')
    contents = contents.replace(remove.get_text(), '')
    for rm in remove2:
        contents = contents.replace(rm.get_text(), '')

    print(title, href, date, contents)