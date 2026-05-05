import os, httpx
from dotenv import load_dotenv

def main():
    load_dotenv()
    token = os.environ["GITHUB_TOKEN"]
    url = "https://api.github.com/repos/kubernetes/kubernetes/pulls"
    headers = {'Authorization': f'Bearer {token}', 'Accept': 'application/vnd.github+json'}
    params = {'state': 'closed', 'per_page': '30', 'sort': 'updated', 'direction': 'desc'}
    response = httpx.get(url, headers=headers, params=params)
    if response.status_code == 200:
        data = response.json()
        for pr in data:
            if pr['merged_at'] is not None:
                print(f"{pr['number']} | {pr['merged_at']} | {pr['title']}")      

if __name__ == "__main__":
    main()
