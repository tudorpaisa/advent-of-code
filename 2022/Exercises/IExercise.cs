namespace AdventOfCode.Exercises;


public interface IExercise
{
    public int GetDay();
    public Result ExecutePart1(string inputFile);
    public Result ExecutePart2(string inputFile);
}
