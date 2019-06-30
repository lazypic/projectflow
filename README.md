# projectflow

Lazypic 프로젝트 데이터 입니다.

### 사용법
프로젝튼 관리는 고유한 키가되는 ID(코드명)을 기준으로 관리합니다.

Project 추가

```bash
$ projectflow -add -id circle -bugdet [총버짓] -status wip
```

Project 수정

```bash
$ projectflow -set -id circle -status backup
```

Project 삭제

```bash
$ sudo projectflow -rm -id [projctname]
```

Project 검색

```bash
$ projectflow -search [검색어]
```

### 인수

#### DB셋팅

- region: 기본값 "ap-northeast-2", AWS 리전명
- profile: 기본값 "lazypic", AWS Credentials profile 이름
- table: 기본값 "projectflow", AWS Dynamodb tablbe 이름

#### 모드
- add: add mode on
- set: set mode on
- rm: rm mode on

#### 속성
- id: 프로젝트 코드명
- budget: 버짓
- startdate: 시작일
- enddate: 마감일
- status: 상태
- searchword: 기본값 "", 검색어

#### 기타
- help: 도움말 출력
- updatedate: 사용자 업데이트 날짜를 임의로 변경시 사용

### AWS DB권한 설정
AWS DB접근 권한을 설정할 계정에 아래 권한을 부여합니다.

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ListAndDescribe",
            "Effect": "Allow",
            "Action": [
                "dynamodb:List*",
                "dynamodb:DescribeReservedCapacity*",
                "dynamodb:DescribeLimits",
                "dynamodb:DescribeTimeToLive"
            ],
            "Resource": "*"
        },
        {
            "Sid": "SpecificTable",
            "Effect": "Allow",
            "Action": [
                "dynamodb:BatchGet*",
                "dynamodb:DescribeStream",
                "dynamodb:DescribeTable",
                "dynamodb:Get*",
                "dynamodb:Query",
                "dynamodb:Scan",
                "dynamodb:BatchWrite*",
                "dynamodb:CreateTable",
                "dynamodb:Delete*",
                "dynamodb:Update*",
                "dynamodb:PutItem"
            ],
            "Resource": "arn:aws:dynamodb:*:*:table/projectflow"
        }
    ]
}
```