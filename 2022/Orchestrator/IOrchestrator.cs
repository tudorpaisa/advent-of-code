namespace AdventOfCode.Orchestrator;

public interface IOrchestrator
{
    void CollectExercises();
    void RunExercise(int day, int part, bool runTest, string inputPath);
    void Execute();
    void SetArgs(string[] args);
}
