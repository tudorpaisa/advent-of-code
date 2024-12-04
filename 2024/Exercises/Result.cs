namespace AdventOfCode.Exercises;

public class Result
{
    public object Answer;
    public bool Completed;

    public Result(object answer, bool completed)
    {
        Answer = answer;
        Completed = completed;
    }

    public Result()
    {
        Answer = -1;
        Completed = false;
    }
}
