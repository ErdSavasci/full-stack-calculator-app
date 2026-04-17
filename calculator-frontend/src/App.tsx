import {useState, useEffect, useCallback} from 'react';
import {FaClipboard, FaFloppyDisk, FaRegTrashCan} from "react-icons/fa6";
import './App.css';
import {FaHourglass} from "react-icons/fa";

interface CalculationLog {
  id: string;
  first: string;
  second: string;
  operation: string;
  result: string;
  timestamp: string;
}

interface CalculateResponse {
  result: string;
}

// Fallback to localhost:8080 if the env variable isn't set
const API_BASE = import.meta.env.VITE_CALCULATOR_BACKEND_URL || 'http://localhost:8080/api';

function App() {
  const [displayValue, setDisplayValue] = useState<string>('0');
  const [previousValue, setPreviousValue] = useState<string | null>(null);
  const [operator, setOperator] = useState<string | null>(null);
  const [equation, setEquation] = useState<string>('');
  const [history, setHistory] = useState<CalculationLog[]>([]);
  const [waitingForNewValue, setWaitingForNewValue] = useState<boolean>(false);
  const [isSaving, setIsSaving] = useState(false);
  const [lastSavedSignature, setLastSavedSignature] = useState<string | null>(null);

  const fetchHistory = useCallback(async () => {
    const res = await fetch(`${API_BASE}/history`);
    if (!res.ok) {
      throw new Error('Failed to fetch history');
    }
    const data: CalculationLog[] = await res.json();
    setHistory(data);
  }, []);

  // Fetch history on initial load
  useEffect(() => {
    fetchHistory()
        .then(() => console.debug("History is fetched"))
        .catch((reason) => console.debug("History fetch error! Reason: " + reason));
  }, [fetchHistory]);

  const handleNumber = (num: string) => {
    if (waitingForNewValue) {
      setDisplayValue(num);
      setWaitingForNewValue(false);
    } else {
      const cleanLength = displayValue.replace(/[-.]/g, '').length;
      if (cleanLength >= 7) {
        return; // Silently reject the clicks
      }

      setDisplayValue(prev => prev === '0' ? num : prev + num);
    }
  };

  const handleSave = async () => {
    try {
      // This creates the string: "5-add-10"
      const signatureToSave = `${previousValue}-${operator}-${displayValue}`;

      // Stop if it's currently saving OR if the signature matches the last saved one
      if (isSaving || signatureToSave === lastSavedSignature) return;

      setIsSaving(true);

      // Send the math request to the backend
      const calcRes = await fetch(`${API_BASE}/calculate/${operator}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ first: previousValue, second: displayValue })
      });

      if (!calcRes.ok) {
        console.warn('Calculation failed with status:', calcRes.status);
        setDisplayValue("Error");
        setWaitingForNewValue(true);
        setPreviousValue(null);
        setOperator(null);
        return;
      }

      const calcData: CalculateResponse = await calcRes.json();
      const finalResult = formatHugeResult(calcData.result);

      const res = await fetch(`${API_BASE}/history`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          first: previousValue,
          second: displayValue,
          operation: operator,
          result: finalResult
        }),
      });

      if (!res.ok) {
        setDisplayValue("Error");
        setWaitingForNewValue(true);
        setPreviousValue(null);
        setOperator(null);
        return;
      }

      // Memorize the signature
      setLastSavedSignature(signatureToSave);

      // Clear the display
      handleClear();

      await fetchHistory();
    } catch (error) {
      console.error("Error saving:", error);
    } finally {
      setIsSaving(false);
    }
  };

  // For Binary Operations (+, -, *, /, ^)
  const handleBinaryAction = (op: string) => {
    setPreviousValue(displayValue);
    setOperator(op);
    setEquation(`${displayValue} ${getSymbol(op)}`);
    setWaitingForNewValue(true);
  };

  // For Unary Operations (Square Root, Percentage)
  const handleUnaryAction = async (action: string) => {
    try {
      // This creates the string: "5-add-10"
      const signatureToSave = `${previousValue}-${operator}-${displayValue}`;

      const res = await fetch(`${API_BASE}/calculate/${action}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        // Unary endpoints only require the "first" payload!
        body: JSON.stringify({ first: displayValue })
      });

      if (!res.ok) {
        setDisplayValue("Error");
        setWaitingForNewValue(true);
        setPreviousValue(null);
        setOperator(null);
        return;
      }

      // Memorize the signature
      setLastSavedSignature(signatureToSave);

      const data: CalculateResponse = await res.json();
      const finalResult = formatHugeResult(data.result);

      // Formatting the equation for the top display
      const eqStr = action === 'squareroot' ? `√${displayValue} =` : `${displayValue}% =`;

      setDisplayValue(finalResult);
      setEquation(eqStr);
      setWaitingForNewValue(true);

      // Save to history
      await fetch(`${API_BASE}/history`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          first: displayValue,
          second: "", // Blank for unary
          operation: action,
          result: finalResult
        })
      });

      fetchHistory()
          .then(() => console.debug("History is fetched"))
          .catch((reason) => console.debug("History fetch error! Reason: " + reason));
    } catch (error) {
      console.error(`${action} error:`, error);
      setDisplayValue("Error");
      setWaitingForNewValue(true);
      setPreviousValue(null);
      setOperator(null);
    }
  };

  const handleHistoryClear = async () => {
    try {
      const res = await fetch(`${API_BASE}/history`, {
        method: 'DELETE'
      });

      if (!res.ok) {
        setDisplayValue("Error");
        return;
      }

      fetchHistory()
          .then(() => console.debug("History cleared and fetched"))
          .catch((reason) => console.debug("History fetch error! Reason: " + reason));
    } catch (error) {
      console.error(`History error:`, error);
      setDisplayValue("Error");
    }
  };

  const handleCalculate = async () => {
    if (!previousValue || !operator) return;

    try {
      // This creates the string: "5-add-10"
      const signatureToSave = `${previousValue}-${operator}-${displayValue}`;

      // Send the math request to the backend
      const calcRes = await fetch(`${API_BASE}/calculate/${operator}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ first: previousValue, second: displayValue })
      });

      if (!calcRes.ok) {
        console.warn('Calculation failed with status:', calcRes.status);
        setDisplayValue("Error");
        setWaitingForNewValue(true);
        setPreviousValue(null);
        setOperator(null);
        return;
      }

      // Memorize the signature
      setLastSavedSignature(signatureToSave);

      const calcData: CalculateResponse = await calcRes.json();
      const finalResult = formatHugeResult(calcData.result);

      // Update the UI
      setDisplayValue(finalResult);
      setEquation(`${previousValue} ${getSymbol(operator)} ${displayValue} =`);
      setWaitingForNewValue(true);
      setPreviousValue(null);
      setOperator(null);

      // Save the result to the History microservice
      await fetch('http://localhost:8080/api/history', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          first: previousValue,
          second: displayValue,
          operation: operator,
          result: finalResult
        })
      });

      // Refresh the history list on the right pane
      fetchHistory()
          .then(() => console.debug("History is refreshed"))
          .catch((reason) => console.debug("History refresh error! Reason: " + reason));

    } catch (error) {
      console.error("Calculation error:", error);
      setDisplayValue("Error");
      setWaitingForNewValue(true);
      setPreviousValue(null);
      setOperator(null);
    }
  };

  const handleClear = () => {
    setDisplayValue('0');
    setPreviousValue(null);
    setOperator(null);
    setEquation('');
  };

  const handleBackspace = () => {
    setDisplayValue(prev => prev.length > 1 ? (prev.length === 2 && prev.startsWith('-') ? '0' : prev.slice(0, -1)) : '0');
  };

  const toggleSign = () => {
    setDisplayValue(prev => prev.startsWith('-') ? (prev.length > 1 ? prev.substring(1) : '0') : (prev === '0' ? '-' : '-' + prev));
  };

  // Helper to map UI symbols to Go backend routes
  const getSymbol = (op: string) => {
    switch (op) {
      case 'add': return '+';
      case 'subtract': return '-';
      case 'multiply': return '×';
      case 'divide': return '/';
      case 'exponential': return '^';
      case 'squareroot': return '√';
      case 'percentage': return '%';
      default: return '';
    }
  };

  const formatHugeResult = (resultString: string) => {
    const num = parseFloat(resultString);

    // If the number is huge
    if (resultString.length > 7) {
      return num.toExponential(3); // Converts to e+
    }

    return resultString;
  };

  const currentSignature = `${previousValue}-${operator}-${displayValue}`;
  const isDuplicate = currentSignature === lastSavedSignature;

  return (
      <div className="app-container">
        {/* LEFT PANE: Calculator */}
        <div className="calc-section">
          <div className="display">
            <div className="display-equation">{equation}</div>
            <h1 className="display-result">{displayValue}</h1>
            <button className="backspace-btn" onClick={handleBackspace}>←</button>
          </div>

          <div className="keypad">
            {/* ROW 1: Advanced Math & Controls */}
            <button className="btn-blue" onClick={() => handleUnaryAction('percentage')}>%</button>
            <button className="btn-blue" onClick={() => handleUnaryAction('squareroot')}>√</button>
            <button className="btn-blue" onClick={() => handleBinaryAction('exponential')}>xⁿ</button>
            <button className="btn-blue" onClick={handleClear}>AC</button>
            <button className="btn-blue" disabled={isSaving || isDuplicate} onClick={handleSave}>{isSaving ? <FaHourglass size={"1em"} /> : <FaFloppyDisk size={"1em"} />}</button>

            {/* ROW 2 */}
            <button onClick={() => handleNumber('1')}>1</button>
            <button onClick={() => handleNumber('2')}>2</button>
            <button onClick={() => handleNumber('3')}>3</button>
            <button className="btn-blue" onClick={() => handleBinaryAction('divide')}>/</button>
            <button className="btn-blue" onClick={toggleSign}>+/-</button>

            {/* ROW 3 */}
            <button onClick={() => handleNumber('4')}>4</button>
            <button onClick={() => handleNumber('5')}>5</button>
            <button onClick={() => handleNumber('6')}>6</button>
            <button className="btn-blue" onClick={() => handleBinaryAction('multiply')}>X</button>
            <button className="btn-blue btn-tall" onClick={handleCalculate}>=</button>

            {/* ROW 4 */}
            <button onClick={() => handleNumber('7')}>7</button>
            <button onClick={() => handleNumber('8')}>8</button>
            <button onClick={() => handleNumber('9')}>9</button>
            <button className="btn-blue" onClick={() => handleBinaryAction('subtract')}>-</button>

            {/* ROW 5 */}
            <button onClick={() => handleNumber('0')}>0</button>
            <button onClick={() => handleNumber('.')}>.</button>
            <button style={{ visibility: 'hidden' }}> </button> {/* Spacer to keep grid aligned */}
            <button className="btn-blue" onClick={() => handleBinaryAction('add')}>+</button>
          </div>
        </div>

        {/* RIGHT PANE: History */}
        <div className="history-section">
          <div className="history-header">
            <h2>History</h2>
            <button id="history-clear" className="btn-blue" onClick={handleHistoryClear}><FaRegTrashCan size={"0.7em"} /></button>
          </div>

          <div className="history-list">
            {history.map((item) => (
                <div key={item.id} className="history-item">
                  <div className="history-calc">
                <span className="history-equation">
                  {/* Dynamic rendering based on operation type */}
                  {item.operation === 'squareroot'
                      ? `√${item.first}`
                      : item.operation === 'percentage'
                          ? `${item.first}%`
                          : `${item.first} ${getSymbol(item.operation)} ${item.second}`
                  }
                </span>
                    <span className="history-result">{item.result}</span>
                  </div>
                  <button
                      className="history-copy"
                      onClick={() => navigator.clipboard.writeText(item.result)}
                      title="Copy Result"
                  >
                    <FaClipboard size="1em" />
                  </button>
                </div>
            ))}
          </div>

          <div className="history-footer">
            <span className="history-footer-text">(Top 5 entries are stored)</span>
          </div>
        </div>
      </div>
  );
}

export default App;