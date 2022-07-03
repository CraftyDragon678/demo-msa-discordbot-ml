from typing import Union
from fastapi import FastAPI
from transformers import AutoTokenizer, AutoModelForCausalLM

tokenizer = AutoTokenizer.from_pretrained('./save-better', bos_token='</s>', eos_token='</s>', pad_token='<pad>')
model = AutoModelForCausalLM.from_pretrained('./save-better')

app = FastAPI()

@app.get('/')
def health_check():
  return {'status': 'good~'}

def get_answer(text: str):
  result: str = tokenizer.decode(
    model.generate(
      tokenizer.encode('<usr>{str}<sys>', return_tensors='pt'),
      do_sample=True,
      temperature=0.8,
      max_length=64,
    )[0],
  )
  return result.split('<sys>')[1].replace('<pad>', '')

@app.get('/chat')
def chat(message: Union[str, None] = None):
  if message is None:
    return {'message': 'Bad Request'}, 400
  return {'message': get_answer(message)}
