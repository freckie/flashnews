import requests
from bs4 import BeautifulSoup

common_url = 'http://www.etoday.co.kr/main.php/news/flashnews/flash_list?MID=0&varPage={}'
common_url2 = 'http://www.etoday.co.kr/news/flashnews/flash_view?idxno={}'

def get_list(page):
    url = common_url.format(page)
    req = requests.get(url)
    bs = BeautifulSoup(req.text, 'lxml')
    wrapper = bs.find('div', class_='flash_tab_lst')
    items = wrapper.find_all('a')
    return [item.get_text() for item in items]

def get_contents(id):
    url = common_url2.format(id)
    req = requests.get(url)
    bs = BeautifulSoup(req.text, 'lxml')
    wrapper = bs.find('div', class_='articleView')
    return wrapper.get_text()

if __name__ == "__main__":
    li = get_list('1')
    print(get_contentli[0])