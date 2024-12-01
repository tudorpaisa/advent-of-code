using AdventOfCode.Exercises;
using AdventOfCode.Utilities;
using System.Reflection;
using CommandLine;

namespace AdventOfCode.Orchestrator;

public class Orchestrator : IOrchestrator
{
    private CliOptions _cliOptions = new();
    private Dictionary<int, IExercise> _exercises = new();

    public void CollectExercises()
    {
        Dictionary<int, IExercise> exercises = new();

        var instances = from t in Assembly.GetExecutingAssembly().GetTypes()
            where t.IsClass && typeof(IExercise).IsAssignableFrom(t) && t.GetConstructor(Type.EmptyTypes) != null
            select Activator.CreateInstance(t) as IExercise;

        foreach (var i in instances)
        {
            exercises.Add(i.GetDay(), i);
        }

        _exercises = exercises;
    }

    public void RunExercise(int day, int part, bool runTest, string sourcePath)
    {
        Printer.PrintRunDetails(day, part, runTest);

        var exercise = _exercises.GetValueOrDefault(day);
        if (exercise == null)
        {
            throw new ArgumentNullException($"Could not find exercise for day {day}.");
        }

        var queue = new List<Action>();

        switch (part)
        {
            case 1:
                queue.Add(() => Printer.PrintResult(exercise.ExecutePart1(BuildInputFilePath(day, part, runTest, sourcePath))));
                break;
            case 2:
                queue.Add(() => Printer.PrintResult(exercise.ExecutePart2(BuildInputFilePath(day, part, runTest, sourcePath))));
                break;
            default:
                queue.Add(() => Printer.PrintResult(exercise.ExecutePart1(BuildInputFilePath(day, 1, runTest, sourcePath))));
                queue.Add(() => Printer.PrintResult(exercise.ExecutePart2(BuildInputFilePath(day, 2, runTest, sourcePath))));
                break;
        }

        foreach (var i in queue) { i(); }
    }

    private static string BuildInputFilePath(int day, int part, bool runTest, string sourcePath)
    {
        var fileName = Constants.FileNameTemplates.FileNameTemplate
            .Replace("{suffix}", runTest ? Constants.FileNameTemplates.TestFileSuffix : Constants.FileNameTemplates.FileSuffix);
        return Path.Combine(sourcePath, day.ToString(), fileName);
    }

    public void Execute()
    {
        CollectExercises();

        var queue = new List<Action>();

        if (_cliOptions.Day == -1)
        {
            // run all days
            int maxDay = _exercises.Keys.ToList().Max();
            for (int i=1; i<= maxDay; i++)
            {
                if (_exercises.ContainsKey(i))
                {
                    // create a new var where the day number is kept otherwise it passes by reference
                    // and would end executing as `maxday + 1`
                    var _i = i;
                    queue.Add(() => RunExercise(_i, _cliOptions.Part, _cliOptions.RunTest, _cliOptions.InputPath));
                }
            }

        }
        else
        {
            queue.Add(() => RunExercise(_cliOptions.Day, _cliOptions.Part, _cliOptions.RunTest, _cliOptions.InputPath));
        }

        foreach (var i in queue) { i(); }
    }

    public void SetArgs(string[] args)
    {
        Parser.Default.ParseArguments<CliOptions>(args).WithParsed<CliOptions>(o =>
            {
                _cliOptions = o;
            });
    }
}
