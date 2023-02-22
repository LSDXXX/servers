import cfscrape

def test():
	scraper = cfscrape.create_scraper()
	headers = {
		"Content-Type":              "application/json",
		"Authorization":             "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiIzNjk2MTY1MjJAcXEuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWUsImdlb2lwX2NvdW50cnkiOiJLUiJ9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItZ2hHdW9NeUtGNVZtRHdjWjc0OTJLb3JEIn0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJhdXRoMHw2M2U3MDQ2N2Q0ODBjMzc3OWZiZTI2YWYiLCJhdWQiOlsiaHR0cHM6Ly9hcGkub3BlbmFpLmNvbS92MSIsImh0dHBzOi8vb3BlbmFpLm9wZW5haS5hdXRoMGFwcC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNjc2NTEzNjU1LCJleHAiOjE2Nzc3MjMyNTUsImF6cCI6IlRkSkljYmUxNldvVEh0Tjk1bnl5d2g1RTR5T282SXRHIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBtb2RlbC5yZWFkIG1vZGVsLnJlcXVlc3Qgb3JnYW5pemF0aW9uLnJlYWQgb2ZmbGluZV9hY2Nlc3MifQ.ciRGkA3-0c50fV-d-BjwmwE9kdyWNbFg-zXTo0q91lx-gKSrVEYoi2sJGvm6ozp6IIyYYitFKZ2f46LKdPn4uFOtfeR81GJviVru9Ec1Qsql2N0xY6hOCsRX_LN5iA8bzoI-0nZyzJ4VD_a6Zt1oCajn6aYnL1AxZXipALkuuulrGIARJvsXuKo5fWiLAGCBVuDBSxVGXPYCq6o9n6sNrApe89y72Ha1KfHR_EvBZ6Q6X0_coxGFyBjOeUgpimnc8r1sywt65QezEeAQr5rgDGYCkZWuMCwIaasRsTqfrl2XTm1vHJVb4UydeIfuaz6_PxbsJT3ej83eP0-L2RS3KA",
		"Accept":                    "text/event-stream",
		"Referer":                   "https://chat.openai.com/chat",
		"Origin":                    "https://chat.openai.com",
		"User-Agent":                "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36",
		"X-OpenAI-Assistant-App-Id": "",
	}
	data = {
		"action": "variant",
        "messages": [
            {
                "id": "12334",
                "role": "user",
                "content": {"content_type": "text", "parts": ["hello"]},
            }
        ],
        "conversation_id": "123",
        "parent_message_id": "123",
        "model": "text-davinci-002-render"
	}
	res = scraper.get("https://chat.openai.com/backend-api/conversation", data=data, headers=headers)
	print(res.content, res.status_code)

if __name__ == '__main__':
	test()
	