import requests
from groq import Groq
import datetime
from dotenv import dotenv_values
from bs4 import BeautifulSoup
import re
import json
import hashlib
import sys
import time
import PyPDF2
import docx
from dotenv import load_dotenv
import os

env_vars = dotenv_values(".env")
load_dotenv()
GroqAPIKey = os.getenv("GROQ_API_KEY") or env_vars.get("GROQ_API_KEY")
client = Groq(api_key=GroqAPIKey)

def read_document(file_path):
    """Read document from various file formats (txt, pdf, docx)"""
    print("debug message", file=sys.stderr)
    file_extension = file_path.lower().split('.')[-1]
    
    try:
        if file_extension == 'txt':
            with open(file_path, 'r', encoding='utf-8') as f:
                return f.read()
        
        elif file_extension == 'pdf':
            text = ""
            with open(file_path, 'rb') as f:
                pdf_reader = PyPDF2.PdfReader(f)
                for page in pdf_reader.pages:
                    page_text = page.extract_text()
                    if page_text:
                        text += page_text + "\n"
            text = clean_text(text)
            return text
        
        elif file_extension in ['docx', 'doc']:
            doc = docx.Document(file_path)
            text = ""
            for paragraph in doc.paragraphs:
                text += paragraph.text + "\n"
            return text.strip()
        
        else:
            print(f"❌ Unsupported file format: .{file_extension}")
            print("   Supported formats: .txt, .pdf, .docx")
            return None
            
    except Exception as e:
        print(f"❌ Error reading file: {e}")
        return None

def scrape_webpage(url, max_length=4000):
    """Scrape webpage content for market research"""
    try:
        headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        }
        
        response = requests.get(url, headers=headers, timeout=15)
        response.raise_for_status()
        
        soup = BeautifulSoup(response.content, 'html.parser')
        
        for element in soup(['script', 'style', 'nav', 'footer', 'header', 'aside', 'iframe', 'noscript']):
            element.decompose()
        
        main_content = soup.find('article') or soup.find('main') or soup.find('body')
        
        if not main_content:
            return None
        
        text = main_content.get_text(separator='\n', strip=True)
        lines = [line.strip() for line in text.split('\n') if line.strip()]
        text = '\n'.join(lines)
        text = re.sub(r'\n{3,}', '\n\n', text)
        
        if len(text) > max_length:
            text = text[:max_length] + "..."
        
        return text
    except:
        return None

def DuckDuckGoSearch(query):
    """Perform DuckDuckGo search"""
    try:
        url = "https://html.duckduckgo.com/html/"
        params = {"q": query}
        headers = {
            "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
        }
        
        response = requests.post(url, data=params, headers=headers, timeout=10)
        
        if response.status_code == 200:
            text = response.text
            results = []
            parts = text.split('result__a')
            
            for i, part in enumerate(parts[1:6]):
                try:
                    link_start = part.find('href="') + 6
                    link_end = part.find('"', link_start)
                    link = part[link_start:link_end]
                    
                    title_start = part.find('>') + 1
                    title_end = part.find('</a>')
                    title = part[title_start:title_end].strip()
                    
                    snippet_part = parts[i+1] if i+1 < len(parts) else part
                    snippet_start = snippet_part.find('result__snippet')
                    snippet_text = snippet_part[snippet_start:snippet_start+300] if snippet_start > 0 else ""
                    snippet = snippet_text.split('>')[1].split('<')[0].strip() if '>' in snippet_text else "No description"
                    
                    if link and title and not link.startswith('//duckduckgo.com'):
                        results.append({
                            "title": title.replace('&amp;', '&'),
                            "snippet": snippet.replace('&amp;', '&'),
                            "link": link
                        })
                except:
                    continue
            
            return results if results else None
        return None
    except:
        return None

def comprehensive_market_research(startup_idea):
    """Conduct comprehensive market research"""
    print("🔍 Conducting market research...", file=sys.stderr)
    
    research_queries = [
        f"{startup_idea} market size 2024 2025",
        f"{startup_idea} competitors landscape",
        f"{startup_idea} market trends statistics",
        f"{startup_idea} investment funding news",
        f"{startup_idea} market saturation analysis"
    ]
    
    all_research_data = ""
    
    for query in research_queries:
        print(f"  → Searching: {query}")
        results = DuckDuckGoSearch(query)
        
        if results:
            all_research_data += f"\n\n{'='*80}\nRESEARCH QUERY: {query}\n{'='*80}\n"
            
            for i, result in enumerate(results[:3], 1):
                all_research_data += f"\n--- Source {i}: {result['title']} ---\n"
                all_research_data += f"URL: {result['link']}\n"
                all_research_data += f"Summary: {result['snippet']}\n"
                
                content = scrape_webpage(result['link'])
                if content:
                    all_research_data += f"Content:\n{content[:2000]}\n"
                
                all_research_data += "-" * 80 + "\n"
        
        time.sleep(1)
    
    return all_research_data

def analyze_market_viability(startup_doc, research_data):
    """Analyze how well the idea fares in current market"""
    print("\n📊 Analyzing market viability...")
    
    analysis_prompt = f"""You are a senior venture capital analyst with 20 years of experience. Analyze this startup idea with BRUTAL HONESTY.

STARTUP IDEA:
{startup_doc}

CURRENT MARKET DATA (2024-2025):
{research_data[:8000]}

Provide a JSON response with this EXACT structure:
{{
    "market_viability_score": <0-100>,
    "market_timing": "early|perfect|late|too_late",
    "market_size_rating": "tiny|small|medium|large|massive",
    "current_market_conditions": "favorable|neutral|unfavorable|hostile",
    "key_market_insights": ["insight1", "insight2", "insight3"],
    "viability_summary": "2-3 sentences explaining the score"
}}

Be extremely critical. A score of 70+ should be rare. Most ideas are mediocre (40-60 range)."""

    try:
        completion = client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            messages=[
                {"role": "system", "content": "You are a brutally honest VC analyst. Respond ONLY with valid JSON."},
                {"role": "user", "content": analysis_prompt}
            ],
            temperature=0.3,
            max_tokens=1024,
            stream=False
        )
        
        result = completion.choices[0].message.content.strip()
        
        if "```json" in result:
            result = result.split("```json")[1].split("```")[0].strip()
        elif "```" in result:
            result = result.split("```")[1].split("```")[0].strip()
        
        return safe_json_load(result)
    except Exception as e:
        print(f"Error in viability analysis: {e}")
        return None

def analyze_problem_solution_fit(startup_doc, research_data):
    """Analyze where and how the startup can help"""
    print("\n🎯 Analyzing problem-solution fit...")
    
    fit_prompt = f"""You are a product strategy expert. Analyze where this startup can genuinely help.

STARTUP IDEA:
{startup_doc}

MARKET CONTEXT:
{research_data[:8000]}

Provide a JSON response:
{{
    "primary_pain_points_addressed": ["pain1", "pain2", "pain3"],
    "target_market_segments": ["segment1", "segment2"],
    "unique_value_proposition_strength": <0-100>,
    "problem_urgency": "low|medium|high|critical",
    "solution_differentiation": "none|weak|moderate|strong|exceptional",
    "real_world_applicability": "limited|moderate|broad|universal",
    "impact_areas": ["area1", "area2"],
    "help_analysis": "Detailed 3-4 sentence analysis of how and where this actually helps users"
}}

Be honest about the actual value provided. Don't exaggerate."""

    try:
        completion = client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            messages=[
                {"role": "system", "content": "You are a critical product strategist. Respond ONLY with valid JSON."},
                {"role": "user", "content": fit_prompt}
            ],
            temperature=0.3,
            max_tokens=1024,
            stream=False
        )
        
        result = completion.choices[0].message.content.strip()
        
        if "```json" in result:
            result = result.split("```json")[1].split("```")[0].strip()
        elif "```" in result:
            result = result.split("```")[1].split("```")[0].strip()
        
        return safe_json_load(result)
    except Exception as e:
        print(f"Error in fit analysis: {e}")
        return None

def analyze_market_saturation(startup_doc, research_data):
    """Determine if market is overpopulated"""
    print("\n🌊 Analyzing market saturation...")
    
    saturation_prompt = f"""You are a market research expert. Analyze market saturation with NO SUGAR-COATING.

STARTUP IDEA:
{startup_doc}

COMPETITIVE LANDSCAPE DATA:
{research_data[:8000]}

Provide a JSON response:
{{
    "saturation_level": "empty|low|moderate|high|oversaturated|dying",
    "saturation_score": <0-100>,
    "number_of_direct_competitors": "0-5|5-20|20-50|50-100|100+",
    "number_of_indirect_competitors": "0-5|5-20|20-50|50-100|100+",
    "major_players": ["company1", "company2", "company3"],
    "barriers_to_entry": "very_low|low|moderate|high|very_high",
    "market_consolidation_stage": "emerging|growth|mature|declining",
    "competitive_intensity": "low|moderate|high|cutthroat",
    "whitespace_opportunities": ["opportunity1", "opportunity2"],
    "saturation_analysis": "Detailed 4-5 sentence analysis explaining WHY the market is at this saturation level"
}}

If the market is oversaturated, say it clearly. Don't hold back."""

    try:
        completion = client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            messages=[
                {"role": "system", "content": "You are a no-nonsense market analyst. Respond ONLY with valid JSON."},
                {"role": "user", "content": saturation_prompt}
            ],
            temperature=0.3,
            max_tokens=1024,
            stream=False
        )
        
        result = completion.choices[0].message.content.strip()
        
        if "```json" in result:
            result = result.split("```json")[1].split("```")[0].strip()
        elif "```" in result:
            result = result.split("```")[1].split("```")[0].strip()
        
        return safe_json_load(result)
    except Exception as e:
        print(f"Error in saturation analysis: {e}")
        return None

def generate_recommendations(startup_doc, viability, fit, saturation, research_data):
    """Generate actionable recommendations"""
    print("\n💡 Generating recommendations...")
    
    rec_prompt = f"""You are a startup advisor who has seen thousands of companies succeed and fail. Provide ACTIONABLE, SPECIFIC recommendations.

STARTUP IDEA:
{startup_doc}

ANALYSIS RESULTS:
Market Viability Score: {viability.get('market_viability_score', 'N/A')}/100
Market Timing: {viability.get('market_timing', 'N/A')}
Saturation Level: {saturation.get('saturation_level', 'N/A')}
Competitive Intensity: {saturation.get('competitive_intensity', 'N/A')}
Solution Differentiation: {fit.get('solution_differentiation', 'N/A')}

MARKET CONTEXT:
{research_data[:6000]}

Provide a JSON response:
{{
    "overall_verdict": "kill_it|pivot_hard|proceed_with_caution|promising|strong_potential|exceptional",
    "confidence_level": <0-100>,
    "critical_risks": ["risk1", "risk2", "risk3"],
    "strategic_recommendations": [
        {{"priority": "critical|high|medium", "recommendation": "specific action"}},
        {{"priority": "critical|high|medium", "recommendation": "specific action"}},
        {{"priority": "critical|high|medium", "recommendation": "specific action"}}
    ],
    "pivot_suggestions": ["suggestion1", "suggestion2"],
    "differentiation_strategies": ["strategy1", "strategy2"],
    "go_to_market_advice": ["advice1", "advice2"],
    "funding_feasibility": "very_difficult|difficult|moderate|feasible|highly_feasible",
    "recommended_next_steps": ["step1", "step2", "step3"],
    "red_flags": ["flag1", "flag2"],
    "green_flags": ["flag1", "flag2"],
    "executive_summary": "5-6 sentences of brutally honest final assessment and advice"
}}

Be SPECIFIC. Don't give generic advice like "improve your product". Give ACTIONABLE steps."""

    try:
        completion = client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            messages=[
                {"role": "system", "content": "You are a brutally honest startup advisor. Respond ONLY with valid JSON."},
                {"role": "user", "content": rec_prompt}
            ],
            temperature=0.4,
            max_tokens=2048,
            stream=False
        )
        
        result = completion.choices[0].message.content.strip()
        
        if "```json" in result:
            result = result.split("```json")[1].split("```")[0].strip()
        elif "```" in result:
            result = result.split("```")[1].split("```")[0].strip()
        
        return safe_json_load(result)
    except Exception as e:
        print(f"Error in recommendations: {e}")
        return None

def format_insights_report(viability, fit, saturation, recommendations):
    """Format the complete insights report"""
    report = f"""
{'='*80}
🚀 STARTUP INSIGHTS REPORT
{'='*80}
Generated: {datetime.datetime.now().strftime('%B %d, %Y at %H:%M')}

{'='*80}
📊 MARKET VIABILITY ANALYSIS
{'='*80}
Overall Score: {viability.get('market_viability_score', 'N/A')}/100
Market Timing: {viability.get('market_timing', 'N/A').upper()}
Market Size: {viability.get('market_size_rating', 'N/A').upper()}
Current Conditions: {viability.get('current_market_conditions', 'N/A').upper()}

Summary:
{viability.get('viability_summary', 'N/A')}

Key Market Insights:
"""
    for insight in viability.get('key_market_insights', []):
        report += f"  • {insight}\n"
    
    report += f"""
{'='*80}
🎯 PROBLEM-SOLUTION FIT ANALYSIS
{'='*80}
Value Proposition Strength: {fit.get('unique_value_proposition_strength', 'N/A')}/100
Problem Urgency: {fit.get('problem_urgency', 'N/A').upper()}
Solution Differentiation: {fit.get('solution_differentiation', 'N/A').upper()}
Applicability: {fit.get('real_world_applicability', 'N/A').upper()}

Pain Points Addressed:
"""
    for pain in fit.get('primary_pain_points_addressed', []):
        report += f"  • {pain}\n"
    
    report += f"\nTarget Market Segments:\n"
    for segment in fit.get('target_market_segments', []):
        report += f"  • {segment}\n"
    
    report += f"\nImpact Areas:\n"
    for area in fit.get('impact_areas', []):
        report += f"  • {area}\n"
    
    report += f"\nDetailed Analysis:\n{fit.get('help_analysis', 'N/A')}\n"
    
    report += f"""
{'='*80}
🌊 MARKET SATURATION ANALYSIS
{'='*80}
Saturation Level: {saturation.get('saturation_level', 'N/A').upper()}
Saturation Score: {saturation.get('saturation_score', 'N/A')}/100
Direct Competitors: {saturation.get('number_of_direct_competitors', 'N/A')}
Indirect Competitors: {saturation.get('number_of_indirect_competitors', 'N/A')}
Competitive Intensity: {saturation.get('competitive_intensity', 'N/A').upper()}
Barriers to Entry: {saturation.get('barriers_to_entry', 'N/A').upper()}
Market Stage: {saturation.get('market_consolidation_stage', 'N/A').upper()}

Major Players:
"""
    for player in saturation.get('major_players', []):
        report += f"  • {player}\n"
    
    report += f"\nWhitespace Opportunities:\n"
    for opp in saturation.get('whitespace_opportunities', []):
        report += f"  • {opp}\n"
    
    report += f"\nDetailed Analysis:\n{saturation.get('saturation_analysis', 'N/A')}\n"
    
    report += f"""
{'='*80}
💡 RECOMMENDATIONS & VERDICT
{'='*80}
Overall Verdict: {recommendations.get('overall_verdict', 'N/A').upper().replace('_', ' ')}
Confidence Level: {recommendations.get('confidence_level', 'N/A')}/100
Funding Feasibility: {recommendations.get('funding_feasibility', 'N/A').upper().replace('_', ' ')}

🚨 CRITICAL RISKS:
"""
    for risk in recommendations.get('critical_risks', []):
        report += f"  ⚠️  {risk}\n"
    
    report += f"\n🔴 RED FLAGS:\n"
    for flag in recommendations.get('red_flags', []):
        report += f"  ❌ {flag}\n"
    
    report += f"\n🟢 GREEN FLAGS:\n"
    for flag in recommendations.get('green_flags', []):
        report += f"  ✅ {flag}\n"
    
    report += f"\n📋 STRATEGIC RECOMMENDATIONS:\n"
    for rec in recommendations.get('strategic_recommendations', []):
        priority_emoji = "🔥" if rec['priority'] == 'critical' else "⚡" if rec['priority'] == 'high' else "📌"
        report += f"  {priority_emoji} [{rec['priority'].upper()}] {rec['recommendation']}\n"
    
    report += f"\n🔄 PIVOT SUGGESTIONS:\n"
    for pivot in recommendations.get('pivot_suggestions', []):
        report += f"  • {pivot}\n"
    
    report += f"\n🎯 DIFFERENTIATION STRATEGIES:\n"
    for strategy in recommendations.get('differentiation_strategies', []):
        report += f"  • {strategy}\n"
    
    report += f"\n📈 GO-TO-MARKET ADVICE:\n"
    for advice in recommendations.get('go_to_market_advice', []):
        report += f"  • {advice}\n"
    
    report += f"\n✅ RECOMMENDED NEXT STEPS:\n"
    for i, step in enumerate(recommendations.get('recommended_next_steps', []), 1):
        report += f"  {i}. {step}\n"
    
    report += f"""
{'='*80}
📝 EXECUTIVE SUMMARY
{'='*80}
{recommendations.get('executive_summary', 'N/A')}

{'='*80}
END OF REPORT
{'='*80}
"""
    return report


def clean_text(text: str) -> str:
    if not text:
        return ""

    text = text.encode("utf-8", errors="ignore").decode("utf-8")
    text = re.sub(r'[\x00-\x1f\x7f-\x9f]', '', text)
    return text.strip()


def analyze_startup_idea(startup_document_path):
    """Main function to analyze startup idea"""
    print("\n" + "="*80)
    print("🔬 AI STARTUP INSIGHTS ANALYZER")
    print("="*80)
    
    try:
        startup_doc = read_document(startup_document_path)
        startup_doc = startup_doc[:6000]
        
        if not startup_doc:
            return None
        
        file_extension = startup_document_path.split('.')[-1].upper()
        print(f"✅ Loaded startup document: {startup_document_path} ({file_extension})")
        print(f"📄 Document length: {len(startup_doc)} characters\n")
        
        research_data = comprehensive_market_research(startup_doc[:1000])
        
        viability = analyze_market_viability(startup_doc, research_data)
        if "error" in viability:
            return viability
        if not viability:
            print("❌ Failed to analyze market viability")
            return
        
        fit = analyze_problem_solution_fit(startup_doc, research_data)
        if not fit:
            print("❌ Failed to analyze problem-solution fit")
            return
        
        saturation = analyze_market_saturation(startup_doc, research_data)
        if not saturation:
            print("❌ Failed to analyze market saturation")
            return
        
        recommendations = generate_recommendations(startup_doc, viability, fit, saturation, research_data)
        if not recommendations:
            print("❌ Failed to generate recommendations")
            return
        
        print("\n✅ Analysis complete! Generating report...\n")
        
        report = format_insights_report(viability, fit, saturation, recommendations)
        
        output_filename = f"startup_insights_report_{datetime.datetime.now().strftime('%Y%m%d_%H%M%S')}.txt"
        with open(output_filename, 'w', encoding='utf-8') as f:
            f.write(report)
        
        print(report)
        print(f"\n💾 Report saved to: {output_filename}")
        
        return {
            'viability': viability,
            'fit': fit,
            'saturation': saturation,
            'recommendations': recommendations,
            'report': report
        }
        
    except FileNotFoundError:
        print(f"❌ Error: File '{startup_document_path}' not found")
        return None
    except Exception as e:
        print(f"❌ Error during analysis: {e}")
        return None

def safe_json_load(raw: str):
    if not raw:
        return {"error": "empty_llm_response"}

    # Remove markdown fences
    raw = raw.strip()
    raw = re.sub(r"^```(json)?", "", raw)
    raw = re.sub(r"```$", "", raw)

    # Extract first JSON object
    match = re.search(r"\{.*\}", raw, re.DOTALL)
    if not match:
        return {
            "error": "no_json_found",
            "raw_preview": raw[:500]
        }

    try:
        return json.loads(match.group(0))
    except json.JSONDecodeError as e:
        return {
            "error": "invalid_json",
            "message": str(e),
            "raw_preview": raw[:500]
        }