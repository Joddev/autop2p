# autop2p
P2P 투자 서비스를 제공하는 업체들이 자동분산투자 기능을 제공할 수 없게 되었다.
이에 개인이 자동분산투자를 계속할 수 있도록 만든 서비스이다.

매일 1시이에 한번씩 설정에 따라 자동으로 투자가 실행한다.

***서비스 이용중 겪는 모든 상황에 대해서 책임지지 않습니다.***

## Deploy
1. `conf.yaml`을 원하는 값에 맞게 수정
1. [문서](https://www.serverless.com/framework/docs/providers/aws/guide/credentials/)를 참고 하여 AWS profile 설정 후
1. 배포
    ```bash
    make deploy
    ```

### Conf.yaml
- `settings[]`:
  - `username`: 로그인에 사용되는 ID
  - `password`: 로그인에 사용되는 패스워드
  - `company`: P2P 서비스 업체
    - `Honestfund`: [어니스트펀드](https://www.honestfund.kr/)
    - `Peoplefund`: [피플펀드](https://www.peoplefund.co.kr/)
  - `amount`: 한 상품에 투자하는 금액
  - `periodMin`: 투자하는 상품의 최소 개월 수
  - `periodMax`: 투자하는 상품의 최대 개월 수
  - `rateMin`: 투자하는 상품의 최소 연이율
  - `rateMax`: 투자하는 상품의 최대 연이율
  - `categories`: 투자하는 상품 종류
    - `PF`: 건설자금 상품
    - `CorporateCredit`: 법인신용 상품
    - `PersonalCredit`: 개인신용 상품
    - `MortgageRealEstate`: 부동산담보 상품
    - `UNKNOWN`: 그 외 상품

### 업체별 특이사항
- `Honestfund`
  - 여러회차에 나눠서 모으는 상품의 반복 투자를 하지 않도록 구현
- `Peoplefund`
  - 여러회차에 나눠서 모으는 상품의 반복 투자를 하지 않도록 구현
