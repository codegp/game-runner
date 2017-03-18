service TurnInformer {
  void createBot(required i32 botID),
  void destroyBot(required i32 botID),
  void startTurn(required i32 botID),
  void destroy()
}
